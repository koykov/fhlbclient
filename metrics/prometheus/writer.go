package prometheus

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Writer struct{}

func (Writer) HostStatus(host string, status int) {
	hostCounter.WithLabelValues(host, code2string(status)).Inc()
}

func (Writer) HostThrottle(host string) {
	hostCounter.WithLabelValues(host, "throttle").Inc()
}

var (
	hostCounter *prometheus.CounterVec

	statusIndex map[int]int
	statusBuf   []struct {
		code, text string
	}
)

func init() {
	hostCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fhlbclient_host",
		Help: "Status code counters.",
	}, []string{"host", "status"})
	prometheus.MustRegister(hostCounter)

	statuses := []int{
		http.StatusContinue,
		http.StatusSwitchingProtocols,
		http.StatusProcessing,
		http.StatusEarlyHints,

		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,

		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusUseProxy,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,

		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusProxyAuthRequired,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusLengthRequired,
		http.StatusPreconditionFailed,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestURITooLong,
		http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable,
		http.StatusExpectationFailed,
		http.StatusTeapot,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusLocked,
		http.StatusFailedDependency,
		http.StatusTooEarly,
		http.StatusUpgradeRequired,
		http.StatusPreconditionRequired,
		http.StatusTooManyRequests,
		http.StatusRequestHeaderFieldsTooLarge,
		http.StatusUnavailableForLegalReasons,

		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended,
		http.StatusNetworkAuthenticationRequired,
	}
	statusIndex = make(map[int]int, len(statuses))
	for i := 0; i < len(statuses); i++ {
		code := statuses[i]
		statusBuf = append(statusBuf, struct{ code, text string }{
			code: strconv.Itoa(code),
			text: http.StatusText(code)})
		statusIndex[code] = len(statusBuf) - 1
	}
}

func code2string(code int) string {
	i, ok := statusIndex[code]
	if !ok {
		return "0"
	}
	return statusBuf[i].code
}

func code2text(code int) string {
	i, ok := statusIndex[code]
	if !ok {
		return "0"
	}
	return statusBuf[i].text
}

var _ = code2text
