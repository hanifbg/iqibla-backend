package static

const WAMessageTemplate = `
Halo {{.CustomerName}},

Terima kasih telah berbelanja di iQibla Indonesia!
Pesanan Anda telah berhasil dibuat dengan nomor #{{.OrderNumber}}.

Total Pembayaran: Rp {{.TotalAmount}}

Silakan selesaikan pembayaran untuk mengkonfirmasi pesanan Anda.
Anda dapat melihat detail pesanan dan melakukan pembayaran melalui link berikut:

{{.OrderConfirmationLink}}

Apabila ada pertanyaan lebih lanjut, silakan hubungi kami.

Terima kasih,
Tim iQibla Indonesia
`
