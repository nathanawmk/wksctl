package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/wksctl/pkg/kubernetes/config"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

const validConfig = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNU1ETXlNakF6TlRReU5Gb1hEVEk1TURNeE9UQXpOVFF5TkZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBT0IrCm1yd3pNd3RsZUhJWHprdmF6SXUxb1pPUXY4MHMrYUVJbVlIUzVPM29TUkVwSzI5ZnVZZUczUCtwaEQrdjJIREEKbko2TDEwT2MyaHl3cXRQUjFqb0xuSjNaWjdkWWZMQys1enNvWWFUU1BJWTJTWVF5M25KQ003YmtQM1h6dGpvbwpzMGhmamxEQUZjZ2VxK1QzcmlUVjFRTnZXdDFwdG5NUStYZWVwNTRNVU01S1hHd0NRYXphWk9vQmp5c0lCOTcrCmtQR3BucWFkblNXdndIU3M4S1RjM1RtbCt5bmpnYlIwa25vdm53bGh0ck41VDU0SWRSRTFIb2VLajhma3pRWDUKL3Budlh6KzhJY09ETDFPS3ZsMmhSMzNidUgwT3JVY0N3ZFlqdlNPMWxZYkRCQ1JPSTZuRm9lTWtsSmdRdGZpQwpCMmZGNjNwN0JMSzM5UVNERGMwQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFLbC85UEZodWE0Y010SjBJU0IrbHNDZzFGZkwKTWUxc1FValpWTFFvWjJHY1ZSb25qM1ZYSndTQ3ZKY2d3WGpnZzZMdmc2MWNyTE51Y1pjTW1iOW1EaTk1aDNIbgprV2JiZ2dOL01zM1didjkxTGExVHpyZG5LOWNMem5QK25HVE5OOUxoMTNRWWFnd204V0JBSEdDbldDZjdIT0tiCkN4RCtIdGZoVFNLZVE3T2FZdk5VZ2MrQ2ZHd0UvWUkzZ05OSjBvOXB1UkpTQWw3bjVNU05lRE9BVGV4L1RxM08KM3BUVWl4K1E3eTJZaTI1Tnd2dTkvUXNRaTRGSUcxZzdxcXg3a0tSejRYaUFoWXR2N3F3MHhIVllBREVpUnZ0OQprWllZWWR5ekZ5RnBzUkoyRTk4c045bXVWQmE2Z2tPV1ordU82U0paV3VZMlhwYytPU0VOQ1JQdCtZcz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    server: https://172.17.0.2:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJVmpMRnVLcUx5S293RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB4T1RBek1qSXdNelUwTWpSYUZ3MHlNREF6TWpFd016VTBNamRhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTVhbGhXYzBIVk10eExPNnMKTXkvaG9vOXpicng2VnhXZXdacVMrMHNxMVFXME5jQkFSR280Y00vQnpiVVdUVHlKYmFkWDB6MXY3dENSeUt3RwpnSzgvUFhINXM0b1NUVXBKc0FKQTF4TkFQOXRSSjhzSE9ldjZHMjJzUTJYbk1xZnZyc2dWZEUxVWFscVRyanZqCkZjdDRyT2M0aEw1YTg3UHowWTFPaEVNZ1JwWGJVODRqalFhMFNqbS9ESlptaDFVOHpQYmY4UlpEajBCVlJ5VnEKNFp5UmU5S2FtYjU0aXc0RTRsaUpOOURqT0JmeGNqelNpaXpqZWx6NmlsQXl6UW1vcXZHTUozSE5nOVIyYzlONApiM3VLblpPSkRjQ3lnUGF4Tzhma3NualkydEo1Q1JNd0NWeE5TNHVpcTV6SmxXMHRhOHlkbDRyOEdic3dxNEJEClZ6c2Yxd0lEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFOUnpNUDAraENOZGYvN2R0QXdGSU1DVGRhaUl1Z3Q1SlZuNwpUbkNwVkNuc09uN1BRbHBndWtNOWpYWjl2M0ZBdEVMNDg5Q3Yyc0NtUTJHVzBaeThGRmlHcFE3dHpJdTdzVytFCkx6VFZiTU9TQWYzNjJ5TjBJSHprcld2UGJmZDdOMmRYQzhTamFkakpNeHFaalRLdG9KTXllTVlXQjBqaDNuK3kKRlZEOHdvSkVRMlNFZW5aa2Roc0V0cWw3N1dtMFFOZ3lpRkRlSWU2ektLQVAvaVJ3MW9QcDk3QmtqeXQ4NEs3cApJTytpeFNjYjd6VkxiWjhPSnRZM2dsMEVFVGZIeG5hYUk0MStDclVHWGtNV1BZd091dzhlcVpGc0pwSU1iS2FtCmVkOFZOSitCNytueERUcFdGcUtrcVMwQTNuZlNiZmthMmxJb1ROYXpIekpDbjUzQzAxUT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBNWFsaFdjMEhWTXR4TE82c015L2hvbzl6YnJ4NlZ4V2V3WnFTKzBzcTFRVzBOY0JBClJHbzRjTS9CemJVV1RUeUpiYWRYMHoxdjd0Q1J5S3dHZ0s4L1BYSDVzNG9TVFVwSnNBSkExeE5BUDl0Uko4c0gKT2V2NkcyMnNRMlhuTXFmdnJzZ1ZkRTFVYWxxVHJqdmpGY3Q0ck9jNGhMNWE4N1B6MFkxT2hFTWdScFhiVTg0agpqUWEwU2ptL0RKWm1oMVU4elBiZjhSWkRqMEJWUnlWcTRaeVJlOUthbWI1NGl3NEU0bGlKTjlEak9CZnhjanpTCmlpemplbHo2aWxBeXpRbW9xdkdNSjNITmc5UjJjOU40YjN1S25aT0pEY0N5Z1BheE84ZmtzbmpZMnRKNUNSTXcKQ1Z4TlM0dWlxNXpKbFcwdGE4eWRsNHI4R2Jzd3E0QkRWenNmMXdJREFRQUJBb0lCQUF5RWN0M21JdVFvUW43awpjMVpHNGRGdWFDZzg5WjRSZTVtcHh5RVRNNzV3bFNYbzJKZmlBam1EMlZoUTZtcERSbXBIbUszV3gyY0l6eWxVCjF2WGtsMW5PQUlJY29HcStCYzRtRVVxbnJmVE5DMXRUNFl6eW82c2pDeVNSUlV5cGdwTFFMUHN0eTlBUUo4UnoKVnlrMDhkcmFyMlhzeWlCR1NwKzlSKzVGaWxqT2I3VWlHTFFOdXpNNmhpeGNUUTVtUXIrZ2pjeWRwOGdUTzJ0RApJMVU2RGdnQWRjWVF3YnVxVTNaYU4ya29KWTk5dzQ3aU5lOVBZc3RqN1N5c1ZUb3ZjUVh1Z1BVamdWa1NVWk1NCkNQa1o3aEUyaGcydThqSVVZcERqbEF1eUtqZGhFSzdFZzN1NGlFak1HOEhTM0ZnK21NYjA4MXQvWmNhcnVRWDYKZEVnMUVXRUNnWUVBNkU2dmtUaG05anFnU3JGVnZWdnJJNCt0aDU4N05QeVZjcmNrb3lRSHJDNXNkcFRsOGtpRApqQmtSU2NMdVFvck5jMnZFUzhINDVoVnBxUzJEVmxMSytsMzVZUEpCZENDc1pMendma0FKNGJCb3pHZ1JjOEhmClhTaS91UlJRNXBSSGc3RmRrZ2V4dHpQZzQ5cjF0T2ZGMmdKL2lpVktBRUNIUUlNSE5la1VRQU1DZ1lFQS9SV2UKQW8xL1p0ZmNiNUFjSnROMFZpcHFmems5cFBQcnRhbWUxcDhZUHVKS2ZpUHBjK0h0YjluSS9DSVJjb3hmTTNHUApadG9yUFIrcXY3Qm1xMGFZZDNLT3dKY3l3RU1TcW1ueFR6OEF4UlQ0Nkg1V3ErVlNET05LVXpzWVZEVWJjQTBXCm8yKy9nMTBEa3hvajZsQ1BkV29KaXF1dGlzQW5vZytSMU1QRlNwMENnWUJBbFgraDgvaE1CRWlEKzRGR3Y4TkQKZzdKT3Zpb0x0UjBuWTF0QUw4Z2lTbFhGTWVncno1VWk0ZVU0aUlVTTR1SHpjTWFGK1V0bFRCYXYvZ05CZ0lzRgp5QktJclZFZEkraEpxVzJDNi9MVFYrUUt6L1BxSnNBZWVqR3pGcjdYRytvMTVwMkk5N0trcUR1aG5VSXFKVFdRClFwbUtvb3RNUHFSYmZ4SUdIdUtPV1FLQmdRQytoVitHSEc4a1JLdzFjQTlCU3ozdy84MWNLUU0zQWtrWFlMR3EKYitvWXJOSFhVOEdTOHltRFlqZmpWdUk3a1dDNW9XdUt5Z0p5NlR2cFFpcUlGWVVCcHNQQVNCSjBtZ21iTUZYdwppa1ZTR0ErcE5qS1pCUEZYc21OcGRMdEQ2UmJXcTRPM1ZaQ2VtNDd0Vm1oakpISmF1WkNsUzhoQkE1YlNjVllmCkRhR2dJUUtCZ1FEQ2NyU25ZSFJZeTlDSVpId3V2alJBNjljODlDeWt0SjBNa3Vqd2FWLys5bUtSN25oK1BuVFcKekFuSmJTVjUwYUZvTFJUNjhhalZ2U1FDSnF1bEVpY1VHcE11MjNma0IyRmozcFNEUmpQWlV6a1ZQTFFqNzhYVApnOGZ1MDNjckFMaEpaQVovcmJiSHNJTFJCUTMwUEtXZWZWRXBBR2tnN3MwaWNqeDRMa2IrdGc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

const validConfigWithPublicIP = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNU1ETXlNakF6TlRReU5Gb1hEVEk1TURNeE9UQXpOVFF5TkZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBT0IrCm1yd3pNd3RsZUhJWHprdmF6SXUxb1pPUXY4MHMrYUVJbVlIUzVPM29TUkVwSzI5ZnVZZUczUCtwaEQrdjJIREEKbko2TDEwT2MyaHl3cXRQUjFqb0xuSjNaWjdkWWZMQys1enNvWWFUU1BJWTJTWVF5M25KQ003YmtQM1h6dGpvbwpzMGhmamxEQUZjZ2VxK1QzcmlUVjFRTnZXdDFwdG5NUStYZWVwNTRNVU01S1hHd0NRYXphWk9vQmp5c0lCOTcrCmtQR3BucWFkblNXdndIU3M4S1RjM1RtbCt5bmpnYlIwa25vdm53bGh0ck41VDU0SWRSRTFIb2VLajhma3pRWDUKL3Budlh6KzhJY09ETDFPS3ZsMmhSMzNidUgwT3JVY0N3ZFlqdlNPMWxZYkRCQ1JPSTZuRm9lTWtsSmdRdGZpQwpCMmZGNjNwN0JMSzM5UVNERGMwQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFLbC85UEZodWE0Y010SjBJU0IrbHNDZzFGZkwKTWUxc1FValpWTFFvWjJHY1ZSb25qM1ZYSndTQ3ZKY2d3WGpnZzZMdmc2MWNyTE51Y1pjTW1iOW1EaTk1aDNIbgprV2JiZ2dOL01zM1didjkxTGExVHpyZG5LOWNMem5QK25HVE5OOUxoMTNRWWFnd204V0JBSEdDbldDZjdIT0tiCkN4RCtIdGZoVFNLZVE3T2FZdk5VZ2MrQ2ZHd0UvWUkzZ05OSjBvOXB1UkpTQWw3bjVNU05lRE9BVGV4L1RxM08KM3BUVWl4K1E3eTJZaTI1Tnd2dTkvUXNRaTRGSUcxZzdxcXg3a0tSejRYaUFoWXR2N3F3MHhIVllBREVpUnZ0OQprWllZWWR5ekZ5RnBzUkoyRTk4c045bXVWQmE2Z2tPV1ordU82U0paV3VZMlhwYytPU0VOQ1JQdCtZcz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    server: https://1.2.3.4:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJVmpMRnVLcUx5S293RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB4T1RBek1qSXdNelUwTWpSYUZ3MHlNREF6TWpFd016VTBNamRhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTVhbGhXYzBIVk10eExPNnMKTXkvaG9vOXpicng2VnhXZXdacVMrMHNxMVFXME5jQkFSR280Y00vQnpiVVdUVHlKYmFkWDB6MXY3dENSeUt3RwpnSzgvUFhINXM0b1NUVXBKc0FKQTF4TkFQOXRSSjhzSE9ldjZHMjJzUTJYbk1xZnZyc2dWZEUxVWFscVRyanZqCkZjdDRyT2M0aEw1YTg3UHowWTFPaEVNZ1JwWGJVODRqalFhMFNqbS9ESlptaDFVOHpQYmY4UlpEajBCVlJ5VnEKNFp5UmU5S2FtYjU0aXc0RTRsaUpOOURqT0JmeGNqelNpaXpqZWx6NmlsQXl6UW1vcXZHTUozSE5nOVIyYzlONApiM3VLblpPSkRjQ3lnUGF4Tzhma3NualkydEo1Q1JNd0NWeE5TNHVpcTV6SmxXMHRhOHlkbDRyOEdic3dxNEJEClZ6c2Yxd0lEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFOUnpNUDAraENOZGYvN2R0QXdGSU1DVGRhaUl1Z3Q1SlZuNwpUbkNwVkNuc09uN1BRbHBndWtNOWpYWjl2M0ZBdEVMNDg5Q3Yyc0NtUTJHVzBaeThGRmlHcFE3dHpJdTdzVytFCkx6VFZiTU9TQWYzNjJ5TjBJSHprcld2UGJmZDdOMmRYQzhTamFkakpNeHFaalRLdG9KTXllTVlXQjBqaDNuK3kKRlZEOHdvSkVRMlNFZW5aa2Roc0V0cWw3N1dtMFFOZ3lpRkRlSWU2ektLQVAvaVJ3MW9QcDk3QmtqeXQ4NEs3cApJTytpeFNjYjd6VkxiWjhPSnRZM2dsMEVFVGZIeG5hYUk0MStDclVHWGtNV1BZd091dzhlcVpGc0pwSU1iS2FtCmVkOFZOSitCNytueERUcFdGcUtrcVMwQTNuZlNiZmthMmxJb1ROYXpIekpDbjUzQzAxUT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBNWFsaFdjMEhWTXR4TE82c015L2hvbzl6YnJ4NlZ4V2V3WnFTKzBzcTFRVzBOY0JBClJHbzRjTS9CemJVV1RUeUpiYWRYMHoxdjd0Q1J5S3dHZ0s4L1BYSDVzNG9TVFVwSnNBSkExeE5BUDl0Uko4c0gKT2V2NkcyMnNRMlhuTXFmdnJzZ1ZkRTFVYWxxVHJqdmpGY3Q0ck9jNGhMNWE4N1B6MFkxT2hFTWdScFhiVTg0agpqUWEwU2ptL0RKWm1oMVU4elBiZjhSWkRqMEJWUnlWcTRaeVJlOUthbWI1NGl3NEU0bGlKTjlEak9CZnhjanpTCmlpemplbHo2aWxBeXpRbW9xdkdNSjNITmc5UjJjOU40YjN1S25aT0pEY0N5Z1BheE84ZmtzbmpZMnRKNUNSTXcKQ1Z4TlM0dWlxNXpKbFcwdGE4eWRsNHI4R2Jzd3E0QkRWenNmMXdJREFRQUJBb0lCQUF5RWN0M21JdVFvUW43awpjMVpHNGRGdWFDZzg5WjRSZTVtcHh5RVRNNzV3bFNYbzJKZmlBam1EMlZoUTZtcERSbXBIbUszV3gyY0l6eWxVCjF2WGtsMW5PQUlJY29HcStCYzRtRVVxbnJmVE5DMXRUNFl6eW82c2pDeVNSUlV5cGdwTFFMUHN0eTlBUUo4UnoKVnlrMDhkcmFyMlhzeWlCR1NwKzlSKzVGaWxqT2I3VWlHTFFOdXpNNmhpeGNUUTVtUXIrZ2pjeWRwOGdUTzJ0RApJMVU2RGdnQWRjWVF3YnVxVTNaYU4ya29KWTk5dzQ3aU5lOVBZc3RqN1N5c1ZUb3ZjUVh1Z1BVamdWa1NVWk1NCkNQa1o3aEUyaGcydThqSVVZcERqbEF1eUtqZGhFSzdFZzN1NGlFak1HOEhTM0ZnK21NYjA4MXQvWmNhcnVRWDYKZEVnMUVXRUNnWUVBNkU2dmtUaG05anFnU3JGVnZWdnJJNCt0aDU4N05QeVZjcmNrb3lRSHJDNXNkcFRsOGtpRApqQmtSU2NMdVFvck5jMnZFUzhINDVoVnBxUzJEVmxMSytsMzVZUEpCZENDc1pMendma0FKNGJCb3pHZ1JjOEhmClhTaS91UlJRNXBSSGc3RmRrZ2V4dHpQZzQ5cjF0T2ZGMmdKL2lpVktBRUNIUUlNSE5la1VRQU1DZ1lFQS9SV2UKQW8xL1p0ZmNiNUFjSnROMFZpcHFmems5cFBQcnRhbWUxcDhZUHVKS2ZpUHBjK0h0YjluSS9DSVJjb3hmTTNHUApadG9yUFIrcXY3Qm1xMGFZZDNLT3dKY3l3RU1TcW1ueFR6OEF4UlQ0Nkg1V3ErVlNET05LVXpzWVZEVWJjQTBXCm8yKy9nMTBEa3hvajZsQ1BkV29KaXF1dGlzQW5vZytSMU1QRlNwMENnWUJBbFgraDgvaE1CRWlEKzRGR3Y4TkQKZzdKT3Zpb0x0UjBuWTF0QUw4Z2lTbFhGTWVncno1VWk0ZVU0aUlVTTR1SHpjTWFGK1V0bFRCYXYvZ05CZ0lzRgp5QktJclZFZEkraEpxVzJDNi9MVFYrUUt6L1BxSnNBZWVqR3pGcjdYRytvMTVwMkk5N0trcUR1aG5VSXFKVFdRClFwbUtvb3RNUHFSYmZ4SUdIdUtPV1FLQmdRQytoVitHSEc4a1JLdzFjQTlCU3ozdy84MWNLUU0zQWtrWFlMR3EKYitvWXJOSFhVOEdTOHltRFlqZmpWdUk3a1dDNW9XdUt5Z0p5NlR2cFFpcUlGWVVCcHNQQVNCSjBtZ21iTUZYdwppa1ZTR0ErcE5qS1pCUEZYc21OcGRMdEQ2UmJXcTRPM1ZaQ2VtNDd0Vm1oakpISmF1WkNsUzhoQkE1YlNjVllmCkRhR2dJUUtCZ1FEQ2NyU25ZSFJZeTlDSVpId3V2alJBNjljODlDeWt0SjBNa3Vqd2FWLys5bUtSN25oK1BuVFcKekFuSmJTVjUwYUZvTFJUNjhhalZ2U1FDSnF1bEVpY1VHcE11MjNma0IyRmozcFNEUmpQWlV6a1ZQTFFqNzhYVApnOGZ1MDNjckFMaEpaQVovcmJiSHNJTFJCUTMwUEtXZWZWRXBBR2tnN3MwaWNqeDRMa2IrdGc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

const invalidConfigWithSSHBanner = `System is booting up. See pam_nologin(8)
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNU1ETXlNakF6TlRReU5Gb1hEVEk1TURNeE9UQXpOVFF5TkZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBT0IrCm1yd3pNd3RsZUhJWHprdmF6SXUxb1pPUXY4MHMrYUVJbVlIUzVPM29TUkVwSzI5ZnVZZUczUCtwaEQrdjJIREEKbko2TDEwT2MyaHl3cXRQUjFqb0xuSjNaWjdkWWZMQys1enNvWWFUU1BJWTJTWVF5M25KQ003YmtQM1h6dGpvbwpzMGhmamxEQUZjZ2VxK1QzcmlUVjFRTnZXdDFwdG5NUStYZWVwNTRNVU01S1hHd0NRYXphWk9vQmp5c0lCOTcrCmtQR3BucWFkblNXdndIU3M4S1RjM1RtbCt5bmpnYlIwa25vdm53bGh0ck41VDU0SWRSRTFIb2VLajhma3pRWDUKL3Budlh6KzhJY09ETDFPS3ZsMmhSMzNidUgwT3JVY0N3ZFlqdlNPMWxZYkRCQ1JPSTZuRm9lTWtsSmdRdGZpQwpCMmZGNjNwN0JMSzM5UVNERGMwQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFLbC85UEZodWE0Y010SjBJU0IrbHNDZzFGZkwKTWUxc1FValpWTFFvWjJHY1ZSb25qM1ZYSndTQ3ZKY2d3WGpnZzZMdmc2MWNyTE51Y1pjTW1iOW1EaTk1aDNIbgprV2JiZ2dOL01zM1didjkxTGExVHpyZG5LOWNMem5QK25HVE5OOUxoMTNRWWFnd204V0JBSEdDbldDZjdIT0tiCkN4RCtIdGZoVFNLZVE3T2FZdk5VZ2MrQ2ZHd0UvWUkzZ05OSjBvOXB1UkpTQWw3bjVNU05lRE9BVGV4L1RxM08KM3BUVWl4K1E3eTJZaTI1Tnd2dTkvUXNRaTRGSUcxZzdxcXg3a0tSejRYaUFoWXR2N3F3MHhIVllBREVpUnZ0OQprWllZWWR5ekZ5RnBzUkoyRTk4c045bXVWQmE2Z2tPV1ordU82U0paV3VZMlhwYytPU0VOQ1JQdCtZcz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    server: https://172.17.0.2:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJVmpMRnVLcUx5S293RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB4T1RBek1qSXdNelUwTWpSYUZ3MHlNREF6TWpFd016VTBNamRhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTVhbGhXYzBIVk10eExPNnMKTXkvaG9vOXpicng2VnhXZXdacVMrMHNxMVFXME5jQkFSR280Y00vQnpiVVdUVHlKYmFkWDB6MXY3dENSeUt3RwpnSzgvUFhINXM0b1NUVXBKc0FKQTF4TkFQOXRSSjhzSE9ldjZHMjJzUTJYbk1xZnZyc2dWZEUxVWFscVRyanZqCkZjdDRyT2M0aEw1YTg3UHowWTFPaEVNZ1JwWGJVODRqalFhMFNqbS9ESlptaDFVOHpQYmY4UlpEajBCVlJ5VnEKNFp5UmU5S2FtYjU0aXc0RTRsaUpOOURqT0JmeGNqelNpaXpqZWx6NmlsQXl6UW1vcXZHTUozSE5nOVIyYzlONApiM3VLblpPSkRjQ3lnUGF4Tzhma3NualkydEo1Q1JNd0NWeE5TNHVpcTV6SmxXMHRhOHlkbDRyOEdic3dxNEJEClZ6c2Yxd0lEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFOUnpNUDAraENOZGYvN2R0QXdGSU1DVGRhaUl1Z3Q1SlZuNwpUbkNwVkNuc09uN1BRbHBndWtNOWpYWjl2M0ZBdEVMNDg5Q3Yyc0NtUTJHVzBaeThGRmlHcFE3dHpJdTdzVytFCkx6VFZiTU9TQWYzNjJ5TjBJSHprcld2UGJmZDdOMmRYQzhTamFkakpNeHFaalRLdG9KTXllTVlXQjBqaDNuK3kKRlZEOHdvSkVRMlNFZW5aa2Roc0V0cWw3N1dtMFFOZ3lpRkRlSWU2ektLQVAvaVJ3MW9QcDk3QmtqeXQ4NEs3cApJTytpeFNjYjd6VkxiWjhPSnRZM2dsMEVFVGZIeG5hYUk0MStDclVHWGtNV1BZd091dzhlcVpGc0pwSU1iS2FtCmVkOFZOSitCNytueERUcFdGcUtrcVMwQTNuZlNiZmthMmxJb1ROYXpIekpDbjUzQzAxUT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBNWFsaFdjMEhWTXR4TE82c015L2hvbzl6YnJ4NlZ4V2V3WnFTKzBzcTFRVzBOY0JBClJHbzRjTS9CemJVV1RUeUpiYWRYMHoxdjd0Q1J5S3dHZ0s4L1BYSDVzNG9TVFVwSnNBSkExeE5BUDl0Uko4c0gKT2V2NkcyMnNRMlhuTXFmdnJzZ1ZkRTFVYWxxVHJqdmpGY3Q0ck9jNGhMNWE4N1B6MFkxT2hFTWdScFhiVTg0agpqUWEwU2ptL0RKWm1oMVU4elBiZjhSWkRqMEJWUnlWcTRaeVJlOUthbWI1NGl3NEU0bGlKTjlEak9CZnhjanpTCmlpemplbHo2aWxBeXpRbW9xdkdNSjNITmc5UjJjOU40YjN1S25aT0pEY0N5Z1BheE84ZmtzbmpZMnRKNUNSTXcKQ1Z4TlM0dWlxNXpKbFcwdGE4eWRsNHI4R2Jzd3E0QkRWenNmMXdJREFRQUJBb0lCQUF5RWN0M21JdVFvUW43awpjMVpHNGRGdWFDZzg5WjRSZTVtcHh5RVRNNzV3bFNYbzJKZmlBam1EMlZoUTZtcERSbXBIbUszV3gyY0l6eWxVCjF2WGtsMW5PQUlJY29HcStCYzRtRVVxbnJmVE5DMXRUNFl6eW82c2pDeVNSUlV5cGdwTFFMUHN0eTlBUUo4UnoKVnlrMDhkcmFyMlhzeWlCR1NwKzlSKzVGaWxqT2I3VWlHTFFOdXpNNmhpeGNUUTVtUXIrZ2pjeWRwOGdUTzJ0RApJMVU2RGdnQWRjWVF3YnVxVTNaYU4ya29KWTk5dzQ3aU5lOVBZc3RqN1N5c1ZUb3ZjUVh1Z1BVamdWa1NVWk1NCkNQa1o3aEUyaGcydThqSVVZcERqbEF1eUtqZGhFSzdFZzN1NGlFak1HOEhTM0ZnK21NYjA4MXQvWmNhcnVRWDYKZEVnMUVXRUNnWUVBNkU2dmtUaG05anFnU3JGVnZWdnJJNCt0aDU4N05QeVZjcmNrb3lRSHJDNXNkcFRsOGtpRApqQmtSU2NMdVFvck5jMnZFUzhINDVoVnBxUzJEVmxMSytsMzVZUEpCZENDc1pMendma0FKNGJCb3pHZ1JjOEhmClhTaS91UlJRNXBSSGc3RmRrZ2V4dHpQZzQ5cjF0T2ZGMmdKL2lpVktBRUNIUUlNSE5la1VRQU1DZ1lFQS9SV2UKQW8xL1p0ZmNiNUFjSnROMFZpcHFmems5cFBQcnRhbWUxcDhZUHVKS2ZpUHBjK0h0YjluSS9DSVJjb3hmTTNHUApadG9yUFIrcXY3Qm1xMGFZZDNLT3dKY3l3RU1TcW1ueFR6OEF4UlQ0Nkg1V3ErVlNET05LVXpzWVZEVWJjQTBXCm8yKy9nMTBEa3hvajZsQ1BkV29KaXF1dGlzQW5vZytSMU1QRlNwMENnWUJBbFgraDgvaE1CRWlEKzRGR3Y4TkQKZzdKT3Zpb0x0UjBuWTF0QUw4Z2lTbFhGTWVncno1VWk0ZVU0aUlVTTR1SHpjTWFGK1V0bFRCYXYvZ05CZ0lzRgp5QktJclZFZEkraEpxVzJDNi9MVFYrUUt6L1BxSnNBZWVqR3pGcjdYRytvMTVwMkk5N0trcUR1aG5VSXFKVFdRClFwbUtvb3RNUHFSYmZ4SUdIdUtPV1FLQmdRQytoVitHSEc4a1JLdzFjQTlCU3ozdy84MWNLUU0zQWtrWFlMR3EKYitvWXJOSFhVOEdTOHltRFlqZmpWdUk3a1dDNW9XdUt5Z0p5NlR2cFFpcUlGWVVCcHNQQVNCSjBtZ21iTUZYdwppa1ZTR0ErcE5qS1pCUEZYc21OcGRMdEQ2UmJXcTRPM1ZaQ2VtNDd0Vm1oakpISmF1WkNsUzhoQkE1YlNjVllmCkRhR2dJUUtCZ1FEQ2NyU25ZSFJZeTlDSVpId3V2alJBNjljODlDeWt0SjBNa3Vqd2FWLys5bUtSN25oK1BuVFcKekFuSmJTVjUwYUZvTFJUNjhhalZ2U1FDSnF1bEVpY1VHcE11MjNma0IyRmozcFNEUmpQWlV6a1ZQTFFqNzhYVApnOGZ1MDNjckFMaEpaQVovcmJiSHNJTFJCUTMwUEtXZWZWRXBBR2tnN3MwaWNqeDRMa2IrdGc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

func TestSanitizeWithValidConfigAndWithPublicIPChange(t *testing.T) {
	actualConfig, err := config.Sanitize(invalidConfigWithSSHBanner, config.Params{
		APIServerExternalEndpoint: "1.2.3.4",
	})
	assert.NoError(t, err)
	assert.Equal(t, validConfigWithPublicIP, actualConfig)
}

func TestSanitizeWithInvalidConfigContainingSSHBannerAndWithPublicIPChange(t *testing.T) {
	actualConfig, err := config.Sanitize(invalidConfigWithSSHBanner, config.Params{
		APIServerExternalEndpoint: "1.2.3.4",
	})
	assert.NoError(t, err)
	assert.Equal(t, validConfigWithPublicIP, actualConfig)
}

func TestWrite(t *testing.T) {
	testDataDir := "./testdata/"
	testDataPath := "./testdata/test_kubeconfig"
	defer os.RemoveAll(testDataDir)
	validConfigObject, err := clientcmd.Load([]byte(validConfig))
	assert.NoError(t, err)
	_, err = config.Write(testDataPath, *validConfigObject, true)
	assert.NoError(t, err)
	loadedConfig, err := clientcmd.LoadFromFile(testDataPath)
	assert.NoError(t, err)
	err = clientcmd.Validate(*loadedConfig)
	assert.NoError(t, err)
}

func TestMerge(t *testing.T) {
	// Create 2 test kubeconfig objects
	validConfigA, err := clientcmd.Load([]byte(validConfig))
	assert.NoError(t, err)
	validConfigB, err := clientcmd.Load([]byte(validConfigWithPublicIP))
	assert.NoError(t, err)

	mergedConfig := config.Merge(validConfigA, validConfigB)
	err = clientcmd.Validate(*mergedConfig)
	assert.NoError(t, err)
}

func TestInvalidExistingConfig(t *testing.T) {
	// Fail if the current kubeconfig, which will be merged with the newly created one, is invalid
	testDataDir := "./testdata/"
	testDataPath := "./testdata/test_kubeconfig"
	defer os.RemoveAll(testDataDir)
	err := os.Mkdir(testDataDir, 0777)
	assert.NoError(t, err)

	err = ioutil.WriteFile(testDataPath, []byte(invalidConfigWithSSHBanner), 0777)
	assert.NoError(t, err)

	validConfig, err := clientcmd.Load([]byte(validConfig))
	assert.NoError(t, err)

	_, err = config.Write(testDataPath, *validConfig, true)
	assert.Errorf(t, err, "Unable to read existing kubeconfig file")
}
