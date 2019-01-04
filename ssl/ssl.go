package ssl

import (
	"net/http"
	"strings"
	Output "../response"
	Mux "github.com/gorilla/mux"
	"net"
	"fmt"
	"os"
	"time"
	"crypto/tls"
	"math"
)

func DomainRouterHandler(w http.ResponseWriter, r *http.Request) {
	RequestParams := Mux.Vars(r)
	Host := strings.ToLower(RequestParams["host"])

	ips, err := net.LookupIP(Host)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}

	Certs := make([]Output.Cert, 0)

	for _, ip := range ips {
		if IsIPv6(ip.String()) {
			continue
		}

		dialer := net.Dialer{Timeout: 10 * time.Second, Deadline: time.Now().Add(10*time.Second + 5*time.Second)}
		connection, err := tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("[%s]:443", ip), &tls.Config{ServerName: Host})

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s", ip, err)
			continue
		}

		// rembember the checked certs based on their Signature
		signatures := make(map[string]struct{})

		// loop to all certs we get there might be multiple chains, as there may be one or more CAs present on the current system, so we have multiple possible chains
		for _, chain := range connection.ConnectionState().VerifiedChains {
			for _, cert := range chain {
				if _, checked := signatures[string(cert.Signature)]; checked {
					continue
				}

				signatures[string(cert.Signature)] = struct{}{}

				// filter out CA certificates
				if cert.IsCA {
					continue
				}

				RemainingDays := int(math.Ceil(cert.NotAfter.Sub(time.Now()).Hours() / 24))

				Certs = append(Certs, Output.Cert{
					IpAddress:  ip,
					LeftDays:   RemainingDays,
					IssuerName: cert.Issuer.CommonName,
					CommonName: cert.Subject.CommonName,
					Serial:     cert.SerialNumber.String(),
					DNSNames:   cert.DNSNames,
				})
			}
		}

		connection.Close()
	}

	LeftDays := 100000;
	for _, crt := range Certs {
		if LeftDays >= 0 && crt.LeftDays < LeftDays {
			LeftDays = crt.LeftDays
		}
	}

	Response := Output.SslData{
		HostName:     Host,
		LeftDays:     LeftDays,
		Certificates: Certs,
	}

	Output.SendResponse(w, Output.SslHostJsonSuccess{Code: http.StatusOK, Data: Response})
}
