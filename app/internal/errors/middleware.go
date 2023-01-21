package custom_error

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept", "application/json")

		var appError *Error

		err := h(w, r)
		if err != nil {
			if errors.As(err, &appError) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				} else if errors.Is(err, ErrNegativeBalance) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrNegativeBalance.Marshal())
					return
				} else if errors.Is(err, ErrScanFail) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(ErrScanFail.Marshal())
					return
				} else if errors.Is(err, ErrExecuteSQL) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrExecuteSQL.Marshal())
					return
				} else if errors.Is(err, ErrTransaction) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(ErrTransaction.Marshal())
					return
				}

				err = err.(*Error)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(ErrNotFound.Marshal())
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
