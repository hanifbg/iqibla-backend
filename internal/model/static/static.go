package static

const WAMessageTemplate = `
Halo {{.CustomerName}},

Terima kasih telah berbelanja di iQibla Indonesia!
Pesanan Anda telah berhasil dibuat dengan nomor #{{.OrderNumber}}.

Total Pembayaran: Rp{{.TotalAmount}}

Silakan selesaikan pembayaran untuk mengkonfirmasi pesanan Anda.
Anda dapat melihat detail pesanan dan melakukan pembayaran melalui link berikut:

{{.OrderConfirmationLink}}

Apabila ada pertanyaan lebih lanjut, silakan hubungi kami.

Terima kasih,
Tim iQibla Indonesia
`

const TelegramTemplate = `
*📦 New Order Confirmation!*

A new order has been paid and confirmed.

*Order Details:*
- *Order No:* ` + "`{{.OrderNumber}}`" + `
- *Customer:* ` + "`{{.CustomerName}}`" + `
- *Email:* ` + "`{{.CustomerEmail}}`" + `
- *Phone:* ` + "`{{.CustomerPhone}}`" + `
- *Total:* ` + "`Rp{{.TotalAmount}}`" + `

*Items:*
{{range .OrderItems}}
- ` + "`{{.ProductName}}`" + ` (` + "`{{.Quantity}}`x `Rp{{.PriceAtPurchase}}`)" + `
{{end}}

*Shipping Address:*
` + "`{{.ShippingAddress}}`" + `
*Shipping Courier:* ` + "`{{.ShippingCourier}}`" + `
*Shipping Service:* ` + "`{{.ShippingService}}`" + `
`
