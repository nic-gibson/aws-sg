# aws-sg

Simple golang program to update an AWS security group setting a rule to allow access to the public IP address of the network running the program.  If you have AWS rules allowing your home network access and a dynamic IP address, this updates it.

The IP address is obtained by calling out to [https://ifconfig.me](https://ifconfig.me).

I run this hourly as a cron job.

## Syntax

```
aws-sg -r REGION -g GROUP-ID -l RULE-ID
```


## How to use

Log into AWS and set up a security group in EC2. Create a rule to allow access to a single IP address. Take note of the rule id and the security group id and use them as parameters to the script.
