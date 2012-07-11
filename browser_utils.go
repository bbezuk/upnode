package upnode

import (
	"github.com/andelf/go-curl/curl"
	"errors"
	"fmt"
)

func DoHttpCheck(task *Task, opts *Options) (result *Result, eval error) {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	result = new(Result)

	//curl options
	easy.Setopt(curl.OPT_URL, task.address)
	easy.Setopt(curl.OPT_NOSIGNAL, 1)
	easy.Setopt(curl.OPT_USERAGENT, "Uptimo.com Website Performance Monitoring web-agent version-0.2 http://uptimo.com")
	easy.Setopt(curl.OPT_FORBID_REUSE, 1)

	//timeouts
	easy.Setopt(curl.OPT_TIMEOUT, 30)
	easy.Setopt(curl.OPT_CONNECTTIMEOUT, 30)

	if opts.follow {
		easy.Setopt(curl.OPT_FOLLOWLOCATION, 1)
		easy.Setopt(curl.OPT_MAXREDIRS, int(opts.redirects))
	}

	if opts.port > 0 {
		easy.Setopt(curl.OPT_PORT, int(opts.port))
	}

	if !opts.withBody {
		easy.Setopt(curl.OPT_NOBODY, 1)
	}

	err := easy.Perform()

	if err != nil {
		switch err {
		case errors.New(fmt.Sprintf("%d", curl.E_COULDNT_RESOLVE_HOST)):
			result.msg = "DNSE" //DNS Error
		case errors.New(fmt.Sprintf("%d", curl.E_COULDNT_CONNECT)):
			result.msg = "CNC" //Could not connect
		case errors.New(fmt.Sprintf("%d", curl.E_HTTP_RETURNED_ERROR)):
			result.msg = "DNF" // http error
			msg, _ := easy.Getinfo(curl.INFO_RESPONSE_CODE)
			result.httpMsg = msg.(int32)
		case errors.New(fmt.Sprintf("%d", curl.E_OPERATION_TIMEDOUT)):
			result.msg = "TO" // timeout
		case errors.New(fmt.Sprintf("%d", curl.E_TOO_MANY_REDIRECTS)):
			result.msg = "TMR" //  too many redirects
		case errors.New(fmt.Sprintf("%d", curl.E_SEND_ERROR)):
			result.msg = "IESEND" // internal error
		case errors.New(fmt.Sprintf("%d", curl.E_RECV_ERROR)):
			result.msg = "IERECV" // internal error
		case errors.New(fmt.Sprintf("%d", curl.E_BAD_CONTENT_ENCODING)):
			result.msg = "BCE" // bad content encoding
		case errors.New(fmt.Sprintf("%d", curl.E_SSL_CONNECT_ERROR)):
			result.msg = "SSL"
		default:
			result.msg = fmt.Sprintf("UKE,%s", err)
		}
	} else {
		msg, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
		result.httpMsg = int32(msg.(int))
		if err != nil {
			eval = errors.New("CGE") //curl general error
			return
		}

		red, err := easy.Getinfo(curl.INFO_REDIRECT_COUNT)
		result.numRedirects = uint16(red.(int))
		if err != nil {
			eval = errors.New("CGE") //curl general error
			return
		}

		if result.httpMsg == 200 || result.httpMsg == 302 {
			result.success = true
			result.msg = "OK"
		} else {
			result.success = false
			result.msg = fmt.Sprint("%d", result.httpMsg)
		}
	}

	result.timeline = *NewTimeline(easy)

	return

}
