package utils

type Alb struct {
	Type                     string  `json:"type"`
	Time                     string  `json:"time"`
	Elb                      string  `json:"elb"`
	Client_ip                string  `json:"client_ip"`
	Client_port              int     `json:"client_port"`
	Target_ip                string  `json:"target_ip"`
	Target_port              int     `json:"target_port"`
	Request_processing_time  float64 `json:"request_processing_time"`
	Target_processing_time   float64 `json:"target_processing_time"`
	Response_processing_time float64 `json:"response_processing_time"`
	Elb_status_code          int     `json:"elb_status_code"`
	Target_status_code       string  `json:"target_status_code"`
	Received_bytes           int     `json:"received_bytes"`
	Sent_bytes               int     `json:"sent_bytes"`
	Request_verb             string  `json:"request_verb"`
	Request_url              string  `json:"request_url"`
	Request_proto            string  `json:"request_proto"`
	User_agent               string  `json:"user_agent"`
	Ssl_cipher               string  `json:"ssl_cipher"`
	Ssl_protocol             string  `json:"ssl_protocol"`
	Target_group_arn         string  `json:"target_group_arn"`
	Trace_id                 string  `json:"trace_id"`
	Domain_name              string  `json:"domain_name"`
	Chosen_cert_arn          string  `json:"chosen_cert_arn"`
	Matched_rule_priority    string  `json:"matched_rule_priority"`
	Request_creation_time    string  `json:"request_creation_time"`
	Actions_executed         string  `json:"actions_executed"`
	Redirect_url             string  `json:"redirect_url"`
	Lambda_error_reason      string  `json:"lambda_error_reason"`
	Target_port_list         string  `json:"target_port_list"`
	Target_status_code_list  string  `json:"target_status_code_lis"`
	Classification           string  `json:"classification"`
	Classification_reason    string  `json:"classification_reason"`
}
