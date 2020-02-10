gosendmail
==========

A tiny sendmail replacement for environments with smarthosts, based on [mhsendmail](https://github.com/mailhog/mhsendmail)


```bash
> go get github.com/squash/gosendmail

> gosendmail test@mailhog.local <<EOF
From: App <app@mlocal>
To: Test <test@mlocal>
Subject: Test message

Some content!
EOF
```

gosendmail looks for default options in /etc/gosendmail.conf. An example is provided.

You can override the from address (for SMTP `MAIL FROM`):

```bash
gosendmail --from="admin@mailhog.local" test@mailhog.local ...
```

Or pass in multiple recipients:

```bash
gosendmail --from="admin@ocal" test@mailhog.local test2@mailhog.local ...
```

Or override the destination SMTP server:

```bash
gosendmail --smtp-addr="localhost:1026" test@local ...
```

To use from php.ini

```
sendmail_path = "/usr/local/bin/gosendmail"
```

### Licence

Original Copyright ©‎ 2015 - 2016, Ian Kent (http://iankent.uk)
(c) 2016 Josh Grebe (https://github.com/squash/)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
