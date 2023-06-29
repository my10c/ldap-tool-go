# Let's Encrypt

## Background
This document explain how to setup Let’s Encrypt, to obtain a free certificate

### Technologies
The following technologies are used

- running Ubuntu 20.04 or newer
- Let’s Encrypted package
- DNS provided by AWS Route53

### Prerequisite
- the email to use to create the account, we will be using ops@co.my10c.com (example)

- generate a password 16 length (no special charters) 
```
< /dev/urandom tr -dc A-Za-z0-9 | head -c${1:-16}
```

- The domain that will be use is co.my10c.com (example)

### AWS

#### Route 53

- Create a public AWS Route53 zone co.my10c.com 
  will be use to authenticate with Let’s Encrypt certs


#### IAM

- Create a policy and called it cerbot 

- Set description to: Policy for cert - Lets Encrypt cert request and renewal
```

{
    "Version": "2012-10-17",
    "Id": "certbot-dns-route53",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZones",
                "route53:GetChange"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect" : "Allow",
            "Action" : [
                "route53:ChangeResourceRecordSets",
                "route53:ListResourceRecordSets"
            ],
            "Resource" : [
                "arn:aws:route53:::hostedzone/XXXXXXX" <----- CHANGE THIS TO THE CORRENT VALUE
            ]
        }
    ]
}
```

#### Allow an Instance to modify the Route 53 (optional)
- Create an AIM role name cert with same description
	Role -> AWS service -> EC2 -> permissions -> certbot -> Tag {Name : certbot}

- Check if the instance that will run the certbot has the role attached
    if not then attach the role, normally this is done by terraform. 

- To attach IAM role to an instance, make sure the aws-cli is configure correctly
```
1. get the instance id
2. aws --region=<aws-region> ec2 associate-iam-instance-profile --instance-id <instance-id> --iam-instance-profile Name=certbot
```

### Install the package
```
snap install core
snap refresh core
snap install --classic certbot
snap set certbot trust-plugin-with-root=ok
snap install certbot-dns-route53
ln -s /snap/bin/certbot /usr/bin/certbot
```

### Configure cerbot
```
certbot certonly --dns-route53 --dns-route53-propagation-seconds 60 -d "*.co.my10c.com" -d "co.my10c.com"
```

output
```

Saving debug log to /var/log/letsencrypt/letsencrypt.log
Plugins selected: Authenticator dns-route53, Installer None
Enter email address (used for urgent renewal and security notices)
 (Enter 'c' to cancel): ops@co.my10c.com

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Please read the Terms of Service at
https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf. You must
agree in order to register with the ACME server. Do you agree?
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: y

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Would you be willing, once your first certificate is successfully issued, to
share your email address with the Electronic Frontier Foundation, a founding
partner of the Let's Encrypt project and the non-profit organization that
develops Certbot? We'd like to send you email about our work encrypting the web,
EFF news, campaigns, and ways to support digital freedom.
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: y
Account registered.
Requesting a certificate for co.my10c.com
Performing the following challenges:
dns-01 challenge for co.my10c.com
Waiting for verification...
Cleaning up challenges
Subscribe to the EFF mailing list (email: ops@co.my10c.com).
We were unable to subscribe you the EFF mailing list because your e-mail address appears to be invalid. You can try again later by visiting https://act.eff.org.

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live.co.my10c.com/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live.co.my10c.com/privkey.pem
   Your certificate will expire on 2021-04-06. To obtain a new or
   tweaked version of this certificate in the future, simply run
   certbot again. To non-interactively renew *all* of your
   certificates, run "certbot renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le

```

#### Certs file
- cert **/etc/letsencrypt/live/co.my10c.com/fullchain.pem**

- ca **/etc/letsencrypt/live/co.my10c.com/chain.pem**

- key **/etc/letsencrypt/live/co.my10c.com/privkey.pem**

### TEST
```
certbot renew --dns-route53 --dry-run
```

### Setup crontab to automatically update the cert (optional)

We need to create cron jobs to update the cert automatically
my the script can be found [here](https://github.com/my10c/ldap-tool-go/blob/main/docs/letsEncrypt/renew_letsencript.sh)
cron entry example (we need to renew every 3 months!)
cron file location: **/etc/cron.d/renew_cert**
script file location: /usr/local/sbin/renew_letsencript.sh
```
# renew the lets encrypt cert
0 6 1 1,3,6,9,12 * 1 * root /usr/local/sbin/renew_letsencript.sh> /tmp/renew_cert_cron.out 2>&1
```

### Manual update cert
```
certbot certonly --manual -d "*.co.my10c.com" -d "co.my10c.com"
```

### The End
Congraculation you should be all set now : 🦄👏
