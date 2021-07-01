package lightning

import (
	"time"
)

var InvoiceListeningTimeout = time.Minute * 150

// ListenForInvoices starts a goroutine that will repeatedly call waitanyinvoice.
// Each payment received will be fed into the client.PaymentHandler function.
// You can change that function in the meantime.
// Or you can set it to nil if you want to stop listening for invoices.
func (ln *Client) ListenForInvoices() {
	go func() {
		for {
			if ln.PaymentHandler == nil {
				return
			}

			res, err := ln.CallWithCustomTimeout(InvoiceListeningTimeout,
				"waitanyinvoice", ln.LastInvoiceIndex)
			if err != nil {
				if _, ok := err.(ErrorTimeout); ok {
					time.Sleep(time.Minute)
				} else {
					time.Sleep(5 * time.Second)
				}
				continue
			}

			index := res.Get("pay_index").Int()
			ln.LastInvoiceIndex = int(index)

			ln.PaymentHandler(res)
		}
	}()
}
