package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/virtual_server"
	"net/http"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMVirtualServerBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	virtualServerName := fmt.Sprintf("acctest_brocadevtm_virtual_server-%d", randomInt)
	resourceName := "brocadevtm_virtual_server.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMVirtualServerCheckDestroy(state, virtualServerName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMVirtualServerNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateOCSPRequired(),
				ExpectError: regexp.MustCompile(`must be one of none, optional, strict`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateOCSPNonce(),
				ExpectError: regexp.MustCompile(`must be one of off, on or strict`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateSSLClientCertHeaders(),
				ExpectError: regexp.MustCompile(`SSL Client Cert Header must be one of all, none or simple`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateServerHonorFallbackSCSV(),
				ExpectError: regexp.MustCompile(`SSL Honor Fallback SCSV must be one of disabled, enabled or use_default`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateCookieDomain(),
				ExpectError: regexp.MustCompile(`Cookie Domain must be one of no_rewrite, set_to_named or set_to_request`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateCookieSecure(),
				ExpectError: regexp.MustCompile(`Cookie Secure must be one of no_modify, set_secure or unset_secure`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateDNSRRSETOrder(),
				ExpectError: regexp.MustCompile(`DNS RRSET Order must be one of cyclic or fixed`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateGZIPCompressLevel(),
				ExpectError: regexp.MustCompile(`Compression level must be a value within 1-9`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateDataFrameSize(),
				ExpectError: regexp.MustCompile(`data_frame_size must be a value within 100-16777206`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateMaxFrameSize(),
				ExpectError: regexp.MustCompile(`max_frame_size must be a value within 16384-16777215`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateETagRewrite(),
				ExpectError: regexp.MustCompile(`ETag Rewrite must be one of wrap, delete, ignore, weaken or wrap`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateMaxBuffer(),
				ExpectError: regexp.MustCompile(`must be within 1024-16777216`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateHeaderTableSize(),
				ExpectError: regexp.MustCompile(`header_table_size must be a value within 4096-1048576`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateSysLogMsgLenLimit(),
				ExpectError: regexp.MustCompile(`msg_len_lemit must be a value within 480-65535`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateChunkOverheadForwarding(),
				ExpectError: regexp.MustCompile(`Chunk Overhead Forwarding must be one of lazy or eager`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateLocationRewrite(),
				ExpectError: regexp.MustCompile(`Location Rewrite must be one of always, if_host_matches or never`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateSIPDangerousRequestsAction(),
				ExpectError: regexp.MustCompile(`Dangerous requests action must be one of forbid, forward or node`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateSIPMode(),
				ExpectError: regexp.MustCompile(`SIP mode must be one of full_gateway, route or sip_gateway`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateSSLRequestClientCert(),
				ExpectError: regexp.MustCompile(`SSL Request Client Cert must be one of dont_request, request or require`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerValidateServerUseSSLSupport(),
				ExpectError: regexp.MustCompile(`must be one of use_default, disabled or enabled`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName),
				ExpectError: regexp.MustCompile(`must be one of client_first, dns, dns_tcp, ftp, http, https, imaps, imapv2, imapv3, imapv4, ldap, ldaps, pop3, pop3s, rtsp, server_first, siptcp, sipudp, smtp, ssl, stream, telnet, udp or udpstreaming`),
			},
			{
				Config: testAccBrocadeVTMVirtualServerCreate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", virtualServerName),
					resource.TestCheckResourceAttr(resourceName, "add_cluster_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_for", "true"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_proto", "true"),
					resource.TestCheckResourceAttr(resourceName, "autodetect_upgrade_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "close_with_rst", "true"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.0", "completionRule1"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.1", "completionRule2"),
					resource.TestCheckResourceAttr(resourceName, "connect_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ftp_force_server_secure", "true"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.0", "testservice"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.1", "testservice2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_any", "true"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.0", "host1"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.1", "host2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.0", "ip1"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.1", "ip2"),
					resource.TestCheckResourceAttr(resourceName, "note", "create acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "pool", "test-pool"),
					resource.TestCheckResourceAttr(resourceName, "port", "50"),
					resource.TestCheckResourceAttr(resourceName, "protection_class", "testProtectionClass"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "dns"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.0", "ruleOne"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.1", "ruleTwo"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.0", "ruleOne"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.1", "ruleTwo"),
					resource.TestCheckResourceAttr(resourceName, "slm_class", "testClass"),
					resource.TestCheckResourceAttr(resourceName, "so_nagle", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl_client_cert_headers", "all"),
					resource.TestCheckResourceAttr(resourceName, "ssl_decrypt", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl_honor_fallback_scsv", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "transparent", "true"),
					resource.TestCheckResourceAttr(resourceName, "error_file", "testErrorFile"),
					resource.TestCheckResourceAttr(resourceName, "expect_starttls", "true"),
					resource.TestCheckResourceAttr(resourceName, "proxy_close", "true"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.name", "profile1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.urls.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.urls.0", "url1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.urls.1", "url2"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.1.name", "profile2"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.1.urls.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.1.urls.0", "url3"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.1.urls.1", "url4"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.1.urls.2", "url5"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.2.name", "profile3"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.2.urls.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.2.urls.0", "url6"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.2.urls.1", "url7"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.2.urls.2", "url8"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.2.urls.3", "url9"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_class", "test"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive", "true"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive_timeout", "500"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_client_buffer", "5000"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_server_buffer", "7000"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_transaction_duration", "5"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.server_first_banner", "testbanner"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.timeout", "4000"),
					resource.TestCheckResourceAttr(resourceName, "cookie.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.domain", "no_rewrite"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.new_domain", "testdomain"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_regex", "testregex"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_replace", "testreplace"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.secure", "no_modify"),
					resource.TestCheckResourceAttr(resourceName, "dns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_client_subnet", "true"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_udpsize", "3000"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.max_udpsize", "4000"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.rrset_order", "cyclic"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.verbose", "true"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.0", "testzone1"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.1", "testzone2"),
					resource.TestCheckResourceAttr(resourceName, "ftp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.data_source_port", "10"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.force_client_secure", "true"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_high", "50"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_low", "5"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.ssl_data", "true"),
					resource.TestCheckResourceAttr(resourceName, "gzip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.compress_level", "5"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.etag_rewrite", "weaken"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.0", "mimetype1"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.1", "mimetype2"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.2", "mimetype3"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.max_size", "4000"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.min_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.no_size", "true"),
					resource.TestCheckResourceAttr(resourceName, "http.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http.0.chunk_overhead_forwarding", "eager"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_regex", "testregex"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_replace", "testlocationreplace"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_rewrite", "never"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_default", "text/html"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_detect", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.connect_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.data_frame_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.header_table_size", "4096"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.0", "header1"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.1", "header2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_default", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.0", "header3"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.1", "header4"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_no_streams", "60"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_open_streams", "120"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_concurrent_streams", "10"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_frame_size", "20000"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_header_padding", "10"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.merge_cookie_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.stream_window_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.0.client_connection_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.format", "logfile/format"),
					resource.TestCheckResourceAttr(resourceName, "log.0.save_all", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.server_connection_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.session_persistence_verbose", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.ssl_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.save_all", "true"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.trace_io", "true"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_high", "50"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_low", "20"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_timeout", "35"),
					resource.TestCheckResourceAttr(resourceName, "sip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.dangerous_requests", "forbid"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.follow_route", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.max_connection_mem", "50"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.mode", "full_gateway"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.rewrite_uri", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_high", "60"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_low", "40"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.timeout_messages", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.transaction_timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.add_http_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.0", "cas1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.1", "cas2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.0", "P256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.1", "P384"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.0", "cas1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.1", "cas2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.2", "cas3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.issuer", "issuerName"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.aia", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.nonce", "strict"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.required", "optional"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.responder_cert", "respondercert"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.signer", "fakesigner"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.url", "fake.url"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.issuer", "issuerName2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.aia", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.nonce", "strict"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.required", "optional"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.responder_cert", "respondercert2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.signer", "fakesigner2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.1.url", "fake2.url"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_max_response_age", "50"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_stapling", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_time_tolerance", "50"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.prefer_sslv3", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.request_client_cert", "request"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.send_close_alerts", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_alt_certificates.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_alt_certificates.0", "testssl001"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_default", "testssl002"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.host", "fakehost1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.certificate", "altcert4"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.0", "altcert1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.1", "altcert2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.2", "altcert3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.1.host", "fakehost2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.1.certificate", "altcert1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.1.alt_certificates.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.1.alt_certificates.0", "altcert4"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.1.alt_certificates.1", "altcert5"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.1.alt_certificates.2", "altcert6"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.2.host", "fakehost3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.2.certificate", "altcert2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.2.alt_certificates.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.2.alt_certificates.0", "altcert7"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.2.alt_certificates.1", "altcert8"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.2.alt_certificates.2", "altcert9"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.signature_algorithms", "RSA_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_ciphers", "SSL_RSA_WITH_AES_128_CBC_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl2", "use_default"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl3", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_1", "use_default"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_2", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.trust_magic", "true"),
					resource.TestCheckResourceAttr(resourceName, "syslog.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.format", "syslog/format"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.ip_end_point", "127.0.0.1:515"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.msg_len_limit", "500"),
					resource.TestCheckResourceAttr(resourceName, "udp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.end_point_persistence", "true"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.port_smp", "true"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.response_datagrams_expected", "-1"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.timeout", "25"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.control_out", "testcontrolout"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.error_page_time", "20"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.max_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.refresh_time", "4"),
				),
			},
			{
				Config: testAccBrocadeVTMVirtualServerUpdate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", virtualServerName),
					resource.TestCheckResourceAttr(resourceName, "add_cluster_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_for", "false"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_proto", "false"),
					resource.TestCheckResourceAttr(resourceName, "autodetect_upgrade_headers", "false"),
					resource.TestCheckResourceAttr(resourceName, "close_with_rst", "false"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.0", "completionRule2"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.1", "completionRule3"),
					resource.TestCheckResourceAttr(resourceName, "connect_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ftp_force_server_secure", "false"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.0", "testservice3"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.1", "testservice4"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_any", "false"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.0", "host3"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.1", "host4"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.0", "ip1"),
					resource.TestCheckResourceAttr(resourceName, "note", "update acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "pool", "test-pool"),
					resource.TestCheckResourceAttr(resourceName, "port", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_class", "testProtectionClassUpdate"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "ftp"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.0", "ruleThree"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.0", "ruleFour"),
					resource.TestCheckResourceAttr(resourceName, "slm_class", "testClassUpdate"),
					resource.TestCheckResourceAttr(resourceName, "so_nagle", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl_client_cert_headers", "simple"),
					resource.TestCheckResourceAttr(resourceName, "ssl_decrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl_honor_fallback_scsv", "use_default"),
					resource.TestCheckResourceAttr(resourceName, "transparent", "false"),
					resource.TestCheckResourceAttr(resourceName, "error_file", "testErrorFileUpdate"),
					resource.TestCheckResourceAttr(resourceName, "expect_starttls", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_close", "false"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.name", "profile1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.urls.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.urls.0", "url4"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.0.urls.1", "url3"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_class", "testUpdate"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive", "false"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive_timeout", "250"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_client_buffer", "5050"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_server_buffer", "7070"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_transaction_duration", "10"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.server_first_banner", "testbannerupdate"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.timeout", "4050"),
					resource.TestCheckResourceAttr(resourceName, "cookie.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.domain", "set_to_named"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.new_domain", "testdomainupdate"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_regex", "testregexupdate"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_replace", "testreplaceupdate"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.secure", "unset_secure"),
					resource.TestCheckResourceAttr(resourceName, "dns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_client_subnet", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_udpsize", "3050"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.max_udpsize", "4050"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.rrset_order", "fixed"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.verbose", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.0", "testzone2"),
					resource.TestCheckResourceAttr(resourceName, "ftp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.data_source_port", "15"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.force_client_secure", "false"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_high", "55"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_low", "7"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.ssl_data", "false"),
					resource.TestCheckResourceAttr(resourceName, "gzip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.compress_level", "8"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.etag_rewrite", "wrap"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.0", "mimetype3"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.max_size", "4050"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.min_size", "50"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.no_size", "false"),
					resource.TestCheckResourceAttr(resourceName, "http.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http.0.chunk_overhead_forwarding", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_regex", "testregexupdate"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_replace", "testlocationreplaceupdate"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_rewrite", "if_host_matches"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_default", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_detect", "false"),
					resource.TestCheckResourceAttr(resourceName, "http2.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.connect_timeout", "75"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.data_frame_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.header_table_size", "4098"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.0", "header3"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.1", "header4"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_default", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.0", "header1"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.1", "header2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_no_streams", "80"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_open_streams", "150"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_concurrent_streams", "15"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_frame_size", "20050"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_header_padding", "13"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.merge_cookie_headers", "false"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.stream_window_size", "201"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.0.client_connection_failures", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.format", "logfile/updateformat"),
					resource.TestCheckResourceAttr(resourceName, "log.0.save_all", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.server_connection_failures", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.session_persistence_verbose", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.ssl_failures", "false"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.save_all", "false"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.trace_io", "false"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_high", "75"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_low", "23"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_timeout", "37"),
					resource.TestCheckResourceAttr(resourceName, "sip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.dangerous_requests", "node"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.follow_route", "false"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.max_connection_mem", "60"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.mode", "sip_gateway"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.rewrite_uri", "false"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_high", "73"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_low", "45"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_timeout", "19"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.timeout_messages", "false"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.transaction_timeout", "23"),
					resource.TestCheckResourceAttr(resourceName, "ssl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.add_http_headers", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.0", "cas2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.0", "P384"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.0", "cas2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.1", "cas3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.issuer", "issuerNameUpdated"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.aia", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.nonce", "off"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.required", "optional"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.responder_cert", "respondercert"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.signer", "fakesigner"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.0.url", "fake.url"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_max_response_age", "55"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_stapling", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_time_tolerance", "55"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_timeout", "25"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.prefer_sslv3", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.request_client_cert", "require"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.send_close_alerts", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_alt_certificates.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_alt_certificates.0", "testssl002"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_default", "testssl001"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.host", "fakehost7"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.certificate", "altcert6"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.0", "altcert5"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.1", "altcert1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_server_cert_host_mapping.0.alt_certificates.2", "altcert3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.signature_algorithms", "ECDSA_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_ciphers", "SSL_RSA_WITH_RC4_128_SHA"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl2", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl3", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_1", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_2", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.trust_magic", "false"),
					resource.TestCheckResourceAttr(resourceName, "syslog.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.format", "syslog/formatupdate"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.ip_end_point", "127.0.0.1:700"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.msg_len_limit", "505"),
					resource.TestCheckResourceAttr(resourceName, "udp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.end_point_persistence", "false"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.port_smp", "false"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.response_datagrams_expected", "50"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.timeout", "33"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.control_out", "testcontroloutupdate"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.error_page_time", "25"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.max_time", "75"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.refresh_time", "9"),
				),
			},
		},
	})
}

func testAccBrocadeVTMVirtualServerCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_virtual_server" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		var vs virtualserver.VirtualServer
		client.WorkWithConfigurationResources()
		err := client.GetByName("virtual_servers", name, &vs)
		if err != nil {
			return fmt.Errorf("Error: Brocade vTM error occurred while retrieving Virtual Server: %s", err)
		}
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("Error: Brocade vTM Virtual Server %s still exists", name)
		}
	}
	return nil
}

func testAccBrocadeVTMVirtualServerExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Virtual Server %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Virtual Server ID not set for %s in resources", name)
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)

		var vs virtualserver.VirtualServer
		client.WorkWithConfigurationResources()
		err := client.GetByName("virtual_servers", name, &vs)
		if err != nil {
			return fmt.Errorf("Brocade vTM Virtual Server - error while retrieving virtual server %s: %s", name, err)
		}
		if client.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("Brocade vTM Virtual Server %s not found on remote vTM", name)
	}
}

func testAccBrocadeVTMVirtualServerNoName() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
pool = "test-pool"
port = 80
}
`
}

func testAccBrocadeVTMVirtualServerValidateOCSPRequired() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
pool = "test-pool"
port = 80

ssl = {
	ocsp_issuers = {
		required = "INVALID"
	}

}
}
`
}

func testAccBrocadeVTMVirtualServerValidateOCSPNonce() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
pool = "test-pool"
port = 80

ssl = {
	ocsp_issuers = {
		nonce = "INVALID"
	}
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateSSLClientCertHeaders() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
pool = "test-pool"
port = 80
ssl_client_cert_headers = "INVALID"
}
`
}

func testAccBrocadeVTMVirtualServerValidateServerHonorFallbackSCSV() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
ssl_honor_fallback_scsv = "INVALID"
}
`
}

func testAccBrocadeVTMVirtualServerValidateCookieDomain() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
cookie = {
	domain = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateCookieSecure() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
cookie = {
	secure = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateDNSRRSETOrder() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
dns = {
	rrset_order = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateGZIPCompressLevel() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
gzip = {
	compress_level = 50
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateDataFrameSize() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
http2 = {
	data_frame_size = 50
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateMaxFrameSize() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
http2 = {
	max_frame_size = 1
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateETagRewrite() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
gzip = {
	etag_rewrite = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateMaxBuffer() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
vs_connection = {
	max_client_buffer = 1
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateHeaderTableSize() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
http2 = {
	header_table_size = 1
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateSysLogMsgLenLimit() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
syslog = {
	msg_len_limit = 1
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateChunkOverheadForwarding() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
http = {
	chunk_overhead_forwarding = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateLocationRewrite() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
http = {
	location_rewrite = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateSIPDangerousRequestsAction() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
sip = {
	dangerous_requests = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateSIPMode() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
sip = {
	mode = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateSSLRequestClientCert() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
ssl = {
	request_client_cert = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerValidateServerUseSSLSupport() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
ssl = {
	ssl_support_ssl2 = "INVALID"
}
}
`
}

func testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
protocol = "SOME_INVALID_PROTOCOL"
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerCreate(virtualServerName string) string {
	return fmt.Sprintf(`



resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	add_cluster_ip = true
	add_x_forwarded_for = true
	add_x_forwarded_proto = true
	autodetect_upgrade_headers = true
	bandwidth_class = "test"
	close_with_rst = true
	completionrules = ["completionRule1","completionRule2"]
	connect_timeout = 50
	enabled = true
	ftp_force_server_secure = true
	glb_services = ["testservice","testservice2"]
	listen_on_any = true
	listen_on_hosts = ["host1","host2"]
	listen_on_traffic_ips = ["ip1","ip2"]
	note = "create acceptance test"
	pool = "test-pool"
	port = 50
	protection_class = "testProtectionClass"
	protocol = "dns"
	request_rules = ["ruleOne","ruleTwo"]
	response_rules = ["ruleOne","ruleTwo"]
	slm_class = "testClass"
	so_nagle = true
	ssl_client_cert_headers = "all"
	ssl_decrypt = true
	ssl_honor_fallback_scsv = "enabled"
	transparent = true
	error_file = "testErrorFile"
	expect_starttls = true
	proxy_close = true

	aptimizer = {
		enabled = true
		profile = [{
			name = "profile1"
			urls = ["url1","url2"]
		},
		{
			name = "profile2"
			urls = ["url3","url4","url5"]
		},
		{
			name = "profile3"
			urls = ["url6","url7","url8","url9"]
		}
		]
	}

	vs_connection = {
		keepalive = true
		keepalive_timeout = 500
		max_client_buffer = 5000
		max_server_buffer = 7000
		max_transaction_duration = 5
		server_first_banner = "testbanner"
		timeout = 4000
	}

	cookie = {
		domain = "no_rewrite"
		new_domain = "testdomain"
		path_regex = "testregex"
		path_replace = "testreplace"
		secure = "no_modify"
	}

	dns = {
		edns_client_subnet = true
		edns_udpsize = 3000
		max_udpsize  = 4000
		rrset_order = "cyclic"
		verbose = true
		zones = ["testzone1","testzone2"]
	}

	ftp = {
		data_source_port = 10
		force_client_secure = true
		port_range_high = 50
		port_range_low = 5
		ssl_data = true
	}

	gzip = {
		compress_level = 5
		enabled = true
		etag_rewrite = "weaken"
		include_mime = ["mimetype1","mimetype2","mimetype3"]
		max_size = 4000
		min_size = 5
		no_size = true
	}

	http = {
		chunk_overhead_forwarding = "eager"
		location_regex = "testregex"
		location_replace = "testlocationreplace"
		location_rewrite = "never"
		mime_default = "text/html"
		mime_detect = true
	}

	http2 = {
		connect_timeout = 50
		data_frame_size = 200
		enabled = true
		header_table_size = 4096
		headers_index_blacklist = ["header1","header2"]
		headers_index_default = true
		headers_index_whitelist = ["header3","header4"]
		idle_timeout_no_streams = 60
		idle_timeout_open_streams = 120
		max_concurrent_streams = 10
		max_frame_size = 20000
		max_header_padding = 10
		merge_cookie_headers = true
		stream_window_size = 200
	}

	log = {
		client_connection_failures = true
		enabled = true
		format = "logfile/format"
		save_all = true
		server_connection_failures = true
		session_persistence_verbose = true
		ssl_failures = true
	}

	recent_connections = {
		enabled = true
		save_all = true
	}

	request_tracing = {
		enabled = true
		trace_io = true
	}

	rtsp = {
		streaming_port_range_high = 50
		streaming_port_range_low = 20
		streaming_timeout = 35
	}

	sip = {
		dangerous_requests = "forbid"
		follow_route = true
		max_connection_mem = 50
		mode = "full_gateway"
		rewrite_uri = true
		streaming_port_range_high = 60
		streaming_port_range_low = 40
		streaming_timeout = 15
		timeout_messages = true
		transaction_timeout = 20
	}

	ssl = {
		add_http_headers = true
		client_cert_cas = ["cas1","cas2"]
		elliptic_curves = ["P256","P384"]
		issued_certs_never_expire = ["cas1","cas2","cas3"]
		ocsp_enable = true

		ocsp_issuers = [
		{
			issuer = "issuerName"
			aia = true
			nonce = "strict"
			required = "optional"
			responder_cert = "respondercert"
			signer = "fakesigner"
			url = "fake.url"
		},
		{
			issuer = "issuerName2"
			aia = true
			nonce = "strict"
			required = "optional"
			responder_cert = "respondercert2"
			signer = "fakesigner2"
			url = "fake2.url"
		},
		]
		ocsp_max_response_age = 50
		ocsp_stapling = true
	    	ocsp_time_tolerance = 50
	    	ocsp_timeout = 20
	    	prefer_sslv3 = true
	    	request_client_cert = "request"
	    	send_close_alerts = true
	    	server_cert_alt_certificates = ["testssl001"]
	    	server_cert_default = "testssl002"

	    	ssl_server_cert_host_mapping = [
		{
		  host = "fakehost1"
		  certificate = "altcert4"
		  alt_certificates = ["altcert1","altcert2","altcert3"]

		},
		{
		  host = "fakehost2"
		  certificate = "altcert1"
		  alt_certificates = ["altcert4","altcert5","altcert6"]
		},
		{
		  host = "fakehost3"
		  certificate = "altcert2"
		  alt_certificates = ["altcert7","altcert8","altcert9"]
		}
		]

		signature_algorithms = "RSA_SHA256"
		ssl_ciphers = "SSL_RSA_WITH_AES_128_CBC_SHA256"
		ssl_support_ssl2 = "use_default"
		ssl_support_ssl3 = "disabled"
		ssl_support_tls1 = "enabled"
		ssl_support_tls1_1 = "use_default"
		ssl_support_tls1_2 = "disabled"
		trust_magic = true
	  }

	  syslog = {
	    enabled = true
	    format = "syslog/format"
	    ip_end_point = "127.0.0.1:515"
	    msg_len_limit = 500
	  }

	  udp = {
	    end_point_persistence = true
	    port_smp = true
	    response_datagrams_expected = -1
	    timeout = 25
          }

	web_cache = {
	    control_out = "testcontrolout"
	    enabled = true
	    error_page_time = 20
	    max_time = 50
	    refresh_time = 4
 	 }
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerUpdate(virtualServerName string) string {
	return fmt.Sprintf(`

resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	add_cluster_ip = false
	add_x_forwarded_for = false
	add_x_forwarded_proto = false
	autodetect_upgrade_headers = false
	bandwidth_class = "testUpdate"
	close_with_rst = false
	completionrules = ["completionRule2","completionRule3"]
	connect_timeout = 100
	enabled = false
	ftp_force_server_secure = false
	glb_services = ["testservice3","testservice4"]
	listen_on_any = false
	listen_on_hosts = ["host3","host4"]
	listen_on_traffic_ips = ["ip1"]
	note = "update acceptance test"
	pool = "test-pool"
	port = 100
	protection_class = "testProtectionClassUpdate"
	protocol = "ftp"
	request_rules = ["ruleThree"]
	response_rules = ["ruleFour"]
	slm_class = "testClassUpdate"
	so_nagle = false
	ssl_client_cert_headers = "simple"
	ssl_decrypt = false
	ssl_honor_fallback_scsv = "use_default"
	transparent = false
	error_file = "testErrorFileUpdate"
	expect_starttls = false
	proxy_close = false

	aptimizer = {
		enabled = true
		profile = [{
			name = "profile1"
			urls = ["url4","url3"]
		}
		]
	}

	vs_connection = {
		keepalive = false
		keepalive_timeout = 250
		max_client_buffer = 5050
		max_server_buffer = 7070
		max_transaction_duration = 10
		server_first_banner = "testbannerupdate"
		timeout = 4050
	}

	cookie = {
		domain = "set_to_named"
		new_domain = "testdomainupdate"
		path_regex = "testregexupdate"
		path_replace = "testreplaceupdate"
		secure = "unset_secure"
	}

	dns = {
		edns_client_subnet = false
		edns_udpsize = 3050
		max_udpsize  = 4050
		rrset_order = "fixed"
		verbose = false
		zones = ["testzone2"]
	}

	ftp = {
		data_source_port = 15
		force_client_secure = false
		port_range_high = 55
		port_range_low = 7
		ssl_data = false
	}

	gzip = {
		compress_level = 8
		enabled = false
		etag_rewrite = "wrap"
		include_mime = ["mimetype3"]
		max_size = 4050
		min_size = 50
		no_size = false
	}

	http = {
		chunk_overhead_forwarding = "lazy"
		location_regex = "testregexupdate"
		location_replace = "testlocationreplaceupdate"
		location_rewrite = "if_host_matches"
		mime_default = "application/json"
		mime_detect = false
	}

	http2 = {
		connect_timeout = 75
		data_frame_size = 100
		enabled = false
		header_table_size = 4098
		headers_index_blacklist = ["header3","header4"]
		headers_index_default = true
		headers_index_whitelist = ["header1","header2"]
		idle_timeout_no_streams = 80
		idle_timeout_open_streams = 150
		max_concurrent_streams = 15
		max_frame_size = 20050
		max_header_padding = 13
		merge_cookie_headers = false
		stream_window_size = 201
	}

	log = {
		client_connection_failures = false
		enabled = false
		format = "logfile/updateformat"
		save_all = false
		server_connection_failures = false
		session_persistence_verbose = false
		ssl_failures = false
	}

	recent_connections = {
		enabled = false
		save_all = false
	}

	request_tracing = {
		enabled = false
		trace_io = false
	}

	rtsp = {
		streaming_port_range_high = 75
		streaming_port_range_low = 23
		streaming_timeout = 37
	}

	sip = {
		dangerous_requests = "node"
		follow_route = false
		max_connection_mem = 60
		mode = "sip_gateway"
		rewrite_uri = false
		streaming_port_range_high = 73
		streaming_port_range_low = 45
		streaming_timeout = 19
		timeout_messages = false
		transaction_timeout = 23
	}

	ssl = {
		add_http_headers = false
		client_cert_cas = ["cas2"]
		elliptic_curves = ["P384"]
		issued_certs_never_expire = ["cas2","cas3"]
		ocsp_enable = false

		ocsp_issuers = [
		{
			issuer = "issuerNameUpdated"
			aia = true
			nonce = "off"
			required = "optional"
			responder_cert = "respondercert"
			signer = "fakesigner"
			url = "fake.url"
		},
		]
		ocsp_max_response_age = 55
		ocsp_stapling = false
	    	ocsp_time_tolerance = 55
	    	ocsp_timeout = 25
	    	prefer_sslv3 = false
	    	request_client_cert = "require"
	    	send_close_alerts = false
	    	server_cert_alt_certificates = ["testssl002"]
	    	server_cert_default = "testssl001"

	    	ssl_server_cert_host_mapping = [
		{
		  host = "fakehost7"
		  certificate = "altcert6"
		  alt_certificates = ["altcert5","altcert1","altcert3"]

		},
		]

		signature_algorithms = "ECDSA_SHA256"
		ssl_ciphers = "SSL_RSA_WITH_RC4_128_SHA"
		ssl_support_ssl2 = "disabled"
		ssl_support_ssl3 = "disabled"
		ssl_support_tls1 = "disabled"
		ssl_support_tls1_1 = "disabled"
		ssl_support_tls1_2 = "disabled"
		trust_magic = false
	  }

	  syslog = {
	    enabled = false
	    format = "syslog/formatupdate"
	    ip_end_point = "127.0.0.1:700"
	    msg_len_limit = 505
	  }

	  udp = {
	    end_point_persistence = false
	    port_smp = false
	    response_datagrams_expected = 50
	    timeout = 33
          }

	web_cache = {
	    control_out = "testcontroloutupdate"
	    enabled = false
	    error_page_time = 25
	    max_time = 75
	    refresh_time = 9
 	 }
}
`, virtualServerName)
}
