go-telegram-widget
==================

[![GoDoc](https://godoc.org/github.com/wesleym/telegramwidget?status.svg)](https://godoc.org/github.com/wesleym/telegramwidget)
[![Build Status](https://github.com/wesleym/telegramwidget/actions/workflows/go.yml/badge.svg)](https://github.com/wesleym/telegramwidget/actions/workflows/go.yml)

This is not an officially supported Google product.

Deprecation
-----------

The [Telegram Login Widget](https://core.telegram.org/widgets/login-legacy) has been deprecated. The updated documentation refers to the login widget as the ["legacy iframe-based JavaScript widget"](https://core.telegram.org/bots/telegram-login#benefits-of-telegram-login).

go-telegram-widget doesn't provide any relevant functionality for Log In With Telegram. Use a library that verifies the signature of JSON Web Tokens (JWTs) instead.

go-telegram-widget remains available for users of the Telegram Login Widget. However, no further updates to this library are planned.

About go-telegram-widget
------------------------

go-telegram-widget simplifies the process of interacting with the [Telegram
login widget](https://core.telegram.org/widgets/login). This library provides
a data type for expressing a Telegram user. It also provides utilities for
interpreting the results of either the JavaScript callback or the URL redirect
into this Telegram user data structure.

This library also verifies the Telegram user data as it interprets it. The
library calls for interpreting the data also perform verification by comparing
with the HMAC code returned from Telegram.
