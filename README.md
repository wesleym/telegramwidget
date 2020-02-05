go-telegram-widget
==================

[![Build Status](https://travis-ci.com/wesleym/telegramwidget.svg?branch=master)](https://travis-ci.com/wesleym/telegramwidget)

This is not an officially supported Google product.

go-telegram-widget simplifies the process of interacting with the [Telegram
login widget](https://core.telegram.org/widgets/login). This library provides
a data type for expressing a Telegram user. It also provides utilities for
interpreting the results of either the JavaScript callback or the URL redirect
into this Telegram user data structure.

This library also verifies the Telegram user data as it interprets it. The
library calls for interpreting the data also perform verification by comparing
with the HMAC code returned from Telegram.
