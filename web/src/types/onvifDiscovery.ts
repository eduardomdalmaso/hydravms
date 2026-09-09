export interface DiscoveredOnvifCamera {
  id: string
  name: string
  manufacturer: string
  model: string
  ip: string
  port: number
  macAddress: string
  profiles: Array<{
    name: string
    token: string
    resolution: string
    codec: string
    rtspUri: string
  }>
  hasPtz: boolean
  isImported: boolean
}
