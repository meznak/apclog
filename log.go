package main

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
)

const (
	// ApacheCommonLog :   {host} {user-identifier} {auth-user-id} [{datetime}] "{method} {request} HTTP/1.0" {response-code} {bytes} {every}
	ApacheCommonLog = "%s - %s %d [%s] \"%s %s\" %d %d %s"
	// ApacheCombinedLog : {host} {user-identifier} {auth-user-id} [{datetime}] "{method} {request} HTTP/1.0" {response-code} {bytes} "{referrer}" "{agent}" {every}
	ApacheCombinedLog = "%s - %s %d [%s] \"%s %s\" %d %d \"%s\" \"%s\" %s"
	// Bluecoat ProxySG
	BlueCoatLog = "<13>date=%s time=%s proxy_name=\"proxy-%d\" proxy_ip=1.2.3.%d duration=%d src=%s user=%s auth_domain=DOMAIN auth_credential_type=- auth_group=domain\\Domain%%20Users supplier_name=%s supplier_ip=%s supplier_country=Unavailable supplier_failures=- exception_id=- exception_category=- filter_result=OBSERVED category=\"%s\" http_referrer=%s status=200 vendor_action=TCP_RESCAN_HIT http_method=GET http_content_type=application/html uri_scheme=https dest_host=%s uri_port=%d uri_path=%s uri_query=- uri_extension=- http_user_agent=%s dvc=1.2.3.%d bytes_in=%d bytes_out=%d virus_id=- symantec_application_name=\"unavailable\" symantec_application_operation=\"unavailable\" risk_score=%d symantec_transaction_uuid=%s icap_reqmod_header=- icap_respmod_header=\"{ %%22expect_sandbox%%22: false  }\" dst_ip=%s dest_ip=%d src_port=%d dst_cert_error=none cs_ocsp_error=- dst_ocsp_error=- dest_ssl_cipher_strength=high dst_cert_hostname=%s dst_cert_category=\"%s\" dst_cert_risk_score=%d proxy_port=3000​ %s"
)

// NewApacheCommonLog creates a log string with apache common log format
func NewApacheCommonLog(t time.Time, every string) string {
	return fmt.Sprintf(
		ApacheCommonLog,
		RandIPV4(),
		RandUser(),
		gofakeit.Number(0, 1000),
		t.Format(time.RFC3339),
		gofakeit.HTTPMethod(),
		RandResourceURI(),
		gofakeit.StatusCode(),
		gofakeit.Number(0, 30000),
		every,
	)
}

// NewApacheCombinedLog creates a log string with apache combined log format
func NewApacheCombinedLog(t time.Time, every string) string {
	return fmt.Sprintf(
		ApacheCombinedLog,
		RandIPV4(),
		RandUser(),
		gofakeit.Number(100, 999),
		t.Format(time.RFC3339),
		gofakeit.HTTPMethod(),
		RandURI(),
		gofakeit.StatusCode(),
		gofakeit.Number(100, 20000),
		RandURL(),
		gofakeit.UserAgent(),
		every,
	)
}

func NewBluecoatLog(t time.Time, every string) string {
	var proxy_number = gofakeit.Number(0, 5)
	var src_ip = RandIPV4()
	var dest_ip = RandIPV4()
	var supplier = RandIPV4()
	var risk_score = gofakeit.Number(0, 10)

	return fmt.Sprintf(
		BlueCoatLog,
		t.Format("2006-01-02"),
		t.Format("15:04:05"),
		proxy_number,
		proxy_number,
		gofakeit.Number(0, 2500),
		src_ip,
		RandUser(),
		supplier,
		supplier,
		RandCategory(),
		RandURI(),
		dest_ip,
		gofakeit.Number(0, 1024),
		RandURI(),
		RandUserAgent(),
		proxy_number,
		gofakeit.Number(0, 5000),
		gofakeit.Number(0, 5000),
		risk_score,
		gofakeit.UUID(),
		dest_ip,
		gofakeit.Number(0, 1024),
		gofakeit.Number(0, 1024),
		RandURL(),
		RandCategory(),
		risk_score,
		every,
	)
}
