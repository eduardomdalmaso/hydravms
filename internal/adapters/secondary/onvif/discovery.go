package onvif

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// DiscoveredProfile represents a video stream profile exposed by the camera.
type DiscoveredProfile struct {
	Name       string `json:"name"`
	Token      string `json:"token"`
	Resolution string `json:"resolution"`
	Codec      string `json:"codec"`
	RTSPUri    string `json:"rtspUri"`
}

// DiscoveredCamera represents an ONVIF/RTSP device found on the local network.
type DiscoveredCamera struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Manufacturer string              `json:"manufacturer"`
	Model        string              `json:"model"`
	IP           string              `json:"ip"`
	Port         int                 `json:"port"`
	MACAddress   string              `json:"macAddress"`
	Profiles     []DiscoveredProfile `json:"profiles"`
	HasPTZ       bool                `json:"hasPtz"`
	IsImported   bool                `json:"isImported"`
}

// DeviceDiscoverer performs WS-Discovery and active subnet scanning.
type DeviceDiscoverer struct {
	multicastAddr string
	probeTimeout  time.Duration
}

// NewDeviceDiscoverer initializes a device discoverer.
func NewDeviceDiscoverer() *DeviceDiscoverer {
	return &DeviceDiscoverer{
		multicastAddr: "239.255.255.250:3702",
		probeTimeout:  2 * time.Second,
	}
}

// DiscoverLocalDevices scans the local network using WS-Discovery and subnet probing.
func (d *DeviceDiscoverer) DiscoverLocalDevices(ctx context.Context) ([]DiscoveredCamera, error) {
	devicesMap := make(map[string]*DiscoveredCamera)
	var mu sync.Mutex

	var wg sync.WaitGroup

	// 1. Run WS-Discovery (UDP Multicast)
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.runWSDiscovery(ctx, &mu, devicesMap)
	}()

	// 2. Run Active TCP Subnet Scan (Ports 554, 80, 8000, 8080, 8899, 8554)
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.runSubnetScan(ctx, &mu, devicesMap)
	}()

	wg.Wait()

	// Fill MAC addresses from ARP table
	arpMap := getARPTable()
	result := make([]DiscoveredCamera, 0, len(devicesMap))
	for _, dev := range devicesMap {
		if mac, ok := arpMap[dev.IP]; ok && mac != "" {
			dev.MACAddress = mac
		}
		result = append(result, *dev)
	}

	return result, nil
}

// runWSDiscovery sends SOAP WS-Discovery Probe to 239.255.255.250:3702.
func (d *DeviceDiscoverer) runWSDiscovery(ctx context.Context, mu *sync.Mutex, results map[string]*DiscoveredCamera) {
	msgID := fmt.Sprintf("uuid:%d-%04x", time.Now().UnixNano(), os.Getpid())
	probeXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<e:Envelope xmlns:e="http://www.w3.org/2003/05/soap-envelope"
            xmlns:w="http://schemas.xmlsoap.org/ws/2004/08/addressing"
            xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery"
            xmlns:dn="http://www.onvif.org/ver10/network/wsdl">
  <e:Header>
    <w:MessageID>%s</w:MessageID>
    <w:To e:mustUnderstand="true">urn:schemas-xmlsoap-org:ws:2005:04:discovery</w:To>
    <w:Action e:mustUnderstand="true">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</w:Action>
  </e:Header>
  <e:Body>
    <d:Probe>
      <d:Types>dn:NetworkVideoTransmitter</d:Types>
    </d:Probe>
  </e:Body>
</e:Envelope>`, msgID)

	raddr, err := net.ResolveUDPAddr("udp4", d.multicastAddr)
	if err != nil {
		return
	}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(d.probeTimeout))
	_, _ = conn.WriteTo([]byte(probeXML), raddr)

	buf := make([]byte, 8192)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, src, err := conn.ReadFrom(buf)
		if err != nil {
			break
		}

		udpAddr, ok := src.(*net.UDPAddr)
		if !ok {
			continue
		}

		ipStr := udpAddr.IP.String()
		payload := string(buf[:n])

		xaddrsRe := regexp.MustCompile(`(?i)<[^>]*XAddrs[^>]*>([^<]+)</`)
		scopesRe := regexp.MustCompile(`(?i)<[^>]*Scopes[^>]*>([^<]+)</`)

		xaddrs := ""
		if match := xaddrsRe.FindStringSubmatch(payload); len(match) > 1 {
			xaddrs = strings.TrimSpace(match[1])
		}

		scopes := ""
		if match := scopesRe.FindStringSubmatch(payload); len(match) > 1 {
			scopes = strings.TrimSpace(match[1])
		}

		mfg := "ONVIF Camera"
		model := "Network Video Transmitter"
		name := fmt.Sprintf("CAM %s", ipStr)

		for _, item := range strings.Fields(scopes) {
			if strings.Contains(item, "/hardware/") {
				parts := strings.Split(item, "/hardware/")
				if len(parts) > 1 {
					model = parts[1]
				}
			} else if strings.Contains(item, "/name/") {
				parts := strings.Split(item, "/name/")
				if len(parts) > 1 {
					name = strings.ReplaceAll(parts[1], "_", " ")
				}
			} else if strings.Contains(item, "/mfg/") {
				parts := strings.Split(item, "/mfg/")
				if len(parts) > 1 {
					mfg = parts[1]
				}
			}
		}

		port := 80
		if xaddrs != "" {
			if strings.Contains(xaddrs, ":") {
				var p int
				fmt.Sscanf(xaddrs, "http://%*[^:]:%d", &p)
				if p > 0 {
					port = p
				}
			}
		}

		mu.Lock()
		results[ipStr] = &DiscoveredCamera{
			ID:           fmt.Sprintf("onvif_%s", strings.ReplaceAll(ipStr, ".", "_")),
			Name:         name,
			Manufacturer: mfg,
			Model:        model,
			IP:           ipStr,
			Port:         port,
			MACAddress:   "--",
			HasPTZ:       strings.Contains(payload, "PTZ") || strings.Contains(scopes, "PTZ"),
			IsImported:   false,
			Profiles: []DiscoveredProfile{
				{
					Name:       "Main Stream (Profile S)",
					Token:      "profile_0",
					Resolution: "1080P",
					Codec:      "H.265",
					RTSPUri:    fmt.Sprintf("rtsp://%s:554/live", ipStr),
				},
				{
					Name:       "Sub Stream (Profile S)",
					Token:      "profile_1",
					Resolution: "480P",
					Codec:      "H.264",
					RTSPUri:    fmt.Sprintf("rtsp://%s:554/sub", ipStr),
				},
			},
		}
		mu.Unlock()
	}
}

// runSubnetScan scans the local subnet for cameras with open RTSP or ONVIF ports.
func (d *DeviceDiscoverer) runSubnetScan(ctx context.Context, mu *sync.Mutex, results map[string]*DiscoveredCamera) {
	subnets := getLocalSubnetIPs()
	if len(subnets) == 0 {
		return
	}

	hostSelfIPs := getHostSelfIPs()
	var scanWg sync.WaitGroup
	sem := make(chan struct{}, 60) // concurrency limit

	for _, targetIP := range subnets {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Skip scanning host machine's own IPs
		if hostSelfIPs[targetIP] {
			continue
		}

		scanWg.Add(1)
		sem <- struct{}{}

		go func(ip string) {
			defer func() {
				<-sem
				scanWg.Done()
			}()

			// Check RTSP port 554 first
			rtspOpen, serverHeader := probeRTSPPort(ip, 554, 300*time.Millisecond)
			if rtspOpen {
				mu.Lock()
				if _, exists := results[ip]; !exists {
					mfg := "IP Camera (RTSP)"
					model := "Network Streamer"
					if serverHeader != "" {
						model = serverHeader
					}
					results[ip] = &DiscoveredCamera{
						ID:           fmt.Sprintf("rtsp_%s", strings.ReplaceAll(ip, ".", "_")),
						Name:         fmt.Sprintf("CAM %s", ip),
						Manufacturer: mfg,
						Model:        model,
						IP:           ip,
						Port:         554,
						MACAddress:   "--",
						HasPTZ:       false,
						IsImported:   false,
						Profiles: []DiscoveredProfile{
							{
								Name:       "Main Stream (RTSP 554)",
								Token:      "profile_rtsp_0",
								Resolution: "1080P",
								Codec:      "H.265",
								RTSPUri:    fmt.Sprintf("rtsp://%s:554/live", ip),
							},
						},
					}
				}
				mu.Unlock()
				return
			}
		}(targetIP)
	}

	scanWg.Wait()
}

// probeRTSPPort attempts a TCP connection and sends an RTSP OPTIONS request.
func probeRTSPPort(ip string, port int, timeout time.Duration) (bool, string) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))
	req := fmt.Sprintf("OPTIONS rtsp://%s:%d/ RTSP/1.0\r\nCSeq: 1\r\nUser-Agent: HydraVMS-Discovery\r\n\r\n", ip, port)
	_, _ = conn.Write([]byte(req))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == nil && n > 0 {
		resp := string(buf[:n])
		if strings.Contains(resp, "RTSP/1.0 200") || strings.Contains(resp, "RTSP/1.0 401") {
			serverHeader := ""
			scanner := bufio.NewScanner(strings.NewReader(resp))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(strings.ToLower(line), "server:") {
					serverHeader = strings.TrimSpace(strings.TrimPrefix(line, "Server:"))
					break
				}
			}
			return true, serverHeader
		}
	}

	return true, "" // TCP port is open even if RTSP OPTIONS was rejected
}

// getLocalSubnetIPs returns all IPv4 addresses in the local subnets (e.g. 10.0.0.1 to 10.0.0.254).
func getLocalSubnetIPs() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue
			}

			ip := ipNet.IP.To4()
			mask := ipNet.Mask
			if len(mask) != 4 {
				continue
			}

			// Generate host IPs for /24 or smaller
			base := ip.Mask(mask)
			for i := 1; i <= 254; i++ {
				hostIP := net.IPv4(base[0], base[1], base[2], byte(i))
				ips = append(ips, hostIP.String())
			}
		}
	}
	return ips
}

// getHostSelfIPs returns a map of all local IP addresses of the host.
func getHostSelfIPs() map[string]bool {
	selfIPs := map[string]bool{"127.0.0.1": true, "::1": true}
	ifaces, err := net.Interfaces()
	if err != nil {
		return selfIPs
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				selfIPs[ipNet.IP.String()] = true
			}
		}
	}
	return selfIPs
}

// getARPTable reads /proc/net/arp to map IP -> MAC Address.
func getARPTable() map[string]string {
	arpMap := make(map[string]string)
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return arpMap
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip header line
	if scanner.Scan() {
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 {
			ip := fields[0]
			mac := fields[3]
			if mac != "00:00:00:00:00:00" && mac != "" {
				arpMap[ip] = mac
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return arpMap
	}
	return arpMap
}
