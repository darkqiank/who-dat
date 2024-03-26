package lib

import (
	whois "github.com/darkqiank/whois"
	whoisparser "github.com/darkqiank/whois-parser"
)

// GetWhois does a WHOIS lookup for a supplied domain
func GetWhois(domain string) (whoisparser.WhoisInfo, error) {
	c := whois.NewClient()
	raw, err := c.Whois(domain)
	// if err != nil {
	// 	return whoisparser.WhoisInfo{}, err
	// }

	result, err1 := whoisparser.Parse(raw)
	if err1 != nil {
		return whoisparser.WhoisInfo{}, err1
	}

	return result, err
}

// GetChanWhois sends Whois data to a channel
func GetChanWhois(domain string, whoisCh chan<- whoisparser.WhoisInfo, errorCh chan<- error) {
	raw, err := whois.Whois(domain)
	if err != nil {
		// return whoisparser.WhoisInfo{}, err
		whoisCh <- whoisparser.WhoisInfo{}
		errorCh <- err
	}

	result, err := whoisparser.Parse(raw)
	if err != nil {
		// return whoisparser.WhoisInfo{}, err
		whoisCh <- whoisparser.WhoisInfo{}
		errorCh <- err
	}

	whoisCh <- result
}
