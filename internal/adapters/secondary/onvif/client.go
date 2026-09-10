package onvif

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Client handles standard ONVIF WS-Security authentication and SOAP media queries.
type Client struct {
	IP         string
	Port       int
	User       string
	Password   string
	timeOffset time.Duration
	httpClient *http.Client
}

// NewClient creates a new ONVIF SOAP client.
func NewClient(ip string, port int, user, password string) *Client {
	if port <= 0 {
		port = 2020
	}
	return &Client{
		IP:       ip,
		Port:     port,
		User:     user,
		Password: password,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// ONVIFDeviceInfo holds device metadata returned by GetDeviceInformation.
type ONVIFDeviceInfo struct {
	Manufacturer    string
	Model           string
	FirmwareVersion string
	SerialNumber    string
	HardwareID      string
}

// syncTime queries GetSystemDateAndTime to eliminate time-drift authorization failures.
func (c *Client) syncTime(ctx context.Context, endpointURL string) {
	soapXML := `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:tds="http://www.onvif.org/ver10/device/wsdl">
  <s:Body>
    <tds:GetSystemDateAndTime/>
  </s:Body>
</s:Envelope>`

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewBufferString(soapXML))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		return
	}

	bodyStr := string(body)
	year, _ := strconv.Atoi(extractTag(bodyStr, "Year"))
	month, _ := strconv.Atoi(extractTag(bodyStr, "Month"))
	day, _ := strconv.Atoi(extractTag(bodyStr, "Day"))
	hour, _ := strconv.Atoi(extractTag(bodyStr, "Hour"))
	min, _ := strconv.Atoi(extractTag(bodyStr, "Minute"))
	sec, _ := strconv.Atoi(extractTag(bodyStr, "Second"))

	if year > 2020 && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
		camTime := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
		c.timeOffset = camTime.Sub(time.Now().UTC())
	}
}

// BuildWSSEHeader generates standard OASIS WS-Security UsernameToken.
func (c *Client) BuildWSSEHeader(usePlainPassword bool) string {
	if c.User == "" {
		return ""
	}

	created := time.Now().UTC().Add(c.timeOffset).Format("2006-01-02T15:04:05Z")

	if usePlainPassword {
		return fmt.Sprintf(`<s:Header>
    <wsse:Security s:mustUnderstand="1" xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
      <wsse:UsernameToken>
        <wsse:Username>%s</wsse:Username>
        <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">%s</wsse:Password>
        <wsu:Created>%s</wsu:Created>
      </wsse:UsernameToken>
    </wsse:Security>
  </s:Header>`, c.User, c.Password, created)
	}

	nonce := make([]byte, 16)
	_, _ = rand.Read(nonce)
	nonceB64 := base64.StdEncoding.EncodeToString(nonce)

	hasher := sha1.New()
	hasher.Write(nonce)
	hasher.Write([]byte(created))
	hasher.Write([]byte(c.Password))
	digest := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf(`<s:Header>
    <wsse:Security s:mustUnderstand="1" xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
      <wsse:UsernameToken>
        <wsse:Username>%s</wsse:Username>
        <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordDigest">%s</wsse:Password>
        <wsse:Nonce EncodingType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary">%s</wsse:Nonce>
        <wsu:Created>%s</wsu:Created>
      </wsse:UsernameToken>
    </wsse:Security>
  </s:Header>`, c.User, digest, nonceB64, created)
}

// sendSOAPRequest sends a SOAP request to the ONVIF device with retry for Plain/Digest auth.
func (c *Client) sendSOAPRequest(ctx context.Context, endpointURL, bodyXML string) (string, int, error) {
	// Try WS-Security PasswordDigest first
	respStr, code, err := c.doSend(ctx, endpointURL, bodyXML, false)
	if err == nil && (code == 200 || !strings.Contains(respStr, "NotAuthorized")) {
		return respStr, code, nil
	}

	// Retry with PasswordText if digest failed
	if c.User != "" && (code == 400 || code == 401 || strings.Contains(respStr, "NotAuthorized") || strings.Contains(respStr, "Authority failure")) {
		respPlain, codePlain, errPlain := c.doSend(ctx, endpointURL, bodyXML, true)
		if errPlain == nil && codePlain == 200 {
			return respPlain, codePlain, nil
		}
	}

	return respStr, code, err
}

func (c *Client) doSend(ctx context.Context, endpointURL, bodyXML string, usePlainAuth bool) (string, int, error) {
	wsseHeader := c.BuildWSSEHeader(usePlainAuth)
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
            xmlns:tds="http://www.onvif.org/ver10/device/wsdl"
            xmlns:trt="http://www.onvif.org/ver10/media/wsdl"
            xmlns:tt="http://www.onvif.org/ver10/schema">
  %s
  <s:Body>
    %s
  </s:Body>
</s:Envelope>`, wsseHeader, bodyXML)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewBufferString(envelope))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}

	return string(respBytes), resp.StatusCode, nil
}

// AuthenticateAndExtractProfiles authenticates via ONVIF SOAP and returns accurate RTSP URIs.
func (c *Client) AuthenticateAndExtractProfiles(ctx context.Context) ([]DiscoveredProfile, ONVIFDeviceInfo, error) {
	portsToTry := []int{c.Port}
	for _, p := range []int{2020, 80, 8000, 8080, 8899, 5000} {
		if p != c.Port && p != 554 {
			portsToTry = append(portsToTry, p)
		}
	}

	var lastErr error
	var activePort int
	var activeMediaURL string
	var devInfo ONVIFDeviceInfo

	for _, port := range portsToTry {
		deviceServiceURL := fmt.Sprintf("http://%s:%d/onvif/device_service", c.IP, port)

		// Sync camera internal time first
		c.syncTime(ctx, deviceServiceURL)

		// 1. Get Device Information
		getDevInfoXML := `<tds:GetDeviceInformation/>`
		resp, statusCode, err := c.sendSOAPRequest(ctx, deviceServiceURL, getDevInfoXML)
		if err == nil && (statusCode == 200 || (statusCode >= 400 && statusCode < 500)) {
			if strings.Contains(resp, "NotAuthorized") || strings.Contains(resp, "Authority failure") || statusCode == 401 {
				lastErr = fmt.Errorf("ONVIF_UNAUTHORIZED")
				continue
			}

			if statusCode == 200 {
				activePort = port
				devInfo.Manufacturer = extractTag(resp, "Manufacturer")
				devInfo.Model = extractTag(resp, "Model")
				devInfo.FirmwareVersion = extractTag(resp, "FirmwareVersion")
				devInfo.SerialNumber = extractTag(resp, "SerialNumber")
				devInfo.HardwareID = extractTag(resp, "HardwareId")

				// 2. Get Capabilities to locate Media Service
				getCapXML := `<tds:GetCapabilities><tds:Category>Media</tds:Category></tds:GetCapabilities>`
				capResp, _, capErr := c.sendSOAPRequest(ctx, deviceServiceURL, getCapXML)
				if capErr == nil {
					activeMediaURL = extractTag(capResp, "XAddr")
				}
				break
			}
		}
		if err != nil {
			lastErr = err
		}
	}

	if activePort == 0 {
		activePort = c.Port
		if activePort == 554 || activePort == 0 {
			activePort = 2020
		}
	}
	if activeMediaURL == "" {
		activeMediaURL = fmt.Sprintf("http://%s:%d/onvif/media_service", c.IP, activePort)
	}

	// 3. Fetch Media Profiles
	getProfilesXML := `<trt:GetProfiles/>`
	profResp, profStatus, err := c.sendSOAPRequest(ctx, activeMediaURL, getProfilesXML)
	if err != nil || (profStatus >= 400 && profStatus < 500) {
		deviceServiceURL := fmt.Sprintf("http://%s:%d/onvif/device_service", c.IP, activePort)
		profResp, profStatus, err = c.sendSOAPRequest(ctx, deviceServiceURL, getProfilesXML)
	}

	if profStatus == 401 || strings.Contains(profResp, "NotAuthorized") || strings.Contains(profResp, "Authority failure") {
		return nil, devInfo, fmt.Errorf("ONVIF_UNAUTHORIZED")
	}

	if err != nil && lastErr != nil && len(profResp) == 0 {
		return nil, devInfo, lastErr
	}

	// 4. Parse Profiles
	profiles := parseProfilesFromSOAP(profResp)
	if len(profiles) == 0 {
		profiles = append(profiles, DiscoveredProfile{
			Name:       "Main Stream",
			Token:      "Profile_1",
			Resolution: "1080P",
			Codec:      "H.264",
		})
	}

	// 5. Query exact RTSP stream URI for each profile
	for i := range profiles {
		token := profiles[i].Token
		getStreamUriXML := fmt.Sprintf(`<trt:GetStreamUri>
      <trt:StreamSetup>
        <tt:Stream>RTP-Unicast</tt:Stream>
        <tt:Transport>
          <tt:Protocol>RTSP</tt:Protocol>
        </tt:Transport>
      </trt:StreamSetup>
      <trt:ProfileToken>%s</trt:ProfileToken>
    </trt:GetStreamUri>`, token)

		uriResp, uriStatus, uriErr := c.sendSOAPRequest(ctx, activeMediaURL, getStreamUriXML)
		if uriErr == nil && uriStatus == 200 {
			uri := extractTag(uriResp, "Uri")
			if uri != "" {
				profiles[i].RTSPUri = c.InjectCredentialsIntoRTSP(uri)
			}
		}

		if profiles[i].RTSPUri == "" {
			profiles[i].RTSPUri = c.InjectCredentialsIntoRTSP(fmt.Sprintf("rtsp://%s:554/stream%d", c.IP, i+1))
		}
	}

	return profiles, devInfo, nil
}

// InjectCredentialsIntoRTSP embeds user:password into the RTSP URI for downstream ingestion.
func (c *Client) InjectCredentialsIntoRTSP(rtspURI string) string {
	if c.User == "" || strings.Contains(rtspURI, "@") {
		return rtspURI
	}
	authPart := fmt.Sprintf("%s:%s@", c.User, c.Password)
	return strings.Replace(rtspURI, "rtsp://", fmt.Sprintf("rtsp://%s", authPart), 1)
}

func extractTag(xmlContent, tagName string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)<(?:[^:]+:)?%s[^>]*>([^<]+)</`, tagName))
	matches := re.FindStringSubmatch(xmlContent)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func parseProfilesFromSOAP(xmlContent string) []DiscoveredProfile {
	var profiles []DiscoveredProfile

	profileBlockRe := regexp.MustCompile(`(?s)<(?:[^:]+:)?Profiles\s+token="([^"]+)"[^>]*>(.*?)</(?:[^:]+:)?Profiles>`)
	matches := profileBlockRe.FindAllStringSubmatch(xmlContent, -1)

	for _, m := range matches {
		token := m[1]
		block := m[2]

		name := extractTag(block, "Name")
		if name == "" {
			name = fmt.Sprintf("Profile %s", token)
		}

		widthStr := extractTag(block, "Width")
		heightStr := extractTag(block, "Height")
		res := "1080P"
		if widthStr != "" && heightStr != "" {
			w, _ := strconv.Atoi(widthStr)
			h, _ := strconv.Atoi(heightStr)
			if w > 0 && h > 0 {
				res = fmt.Sprintf("%dx%d", w, h)
			}
		}

		encoding := extractTag(block, "Encoding")
		if encoding == "" {
			encoding = "H.264"
		}

		profiles = append(profiles, DiscoveredProfile{
			Name:       name,
			Token:      token,
			Resolution: res,
			Codec:      encoding,
		})
	}

	return profiles
}
