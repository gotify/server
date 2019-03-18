---
name: Questions
about: Have problem setting up gotify server? Feel free to ask here
title: ''
labels: question
assignees: ''

---

<!-- 
Alternative ways to get help:
Official documentation - https://gotify.net/
Community chat - https://matrix.to/#/#gotify:matrix.org
-->

**Have you read the documentation?**
- [ ] Yes, but it does not include related information regarding my question.
- [ ] Yes, but the steps described in the documentation do not work on my machine.
- [ ] Yes, but I am having difficulty understanding it and wants clarification.

**You are setting up gotify in**
- [ ] Docker
- [ ] Linux native platform
- [ ] Windows native platform

<details><summary>Describe your configuration (the presense of reverse proxy, VPN connections, API gateways, etc.)<br><pre>

I have an Apache reverse proxy forwarding requests to gotify. Here I will attach the related parts of my Apache configuration:
<VirtualHost *:443>
...
</VirtualHost >

</pre></details>

<details><summary>Paste your <code>docker run</code> command, <code>docker-compose.yml</code>, <code>config.yml</code>, if applicable (remember to mask sensitive information)</summary><br><pre>

$ docker run -p 80:80 -v /var/gotify/data:/app/data gotify/server

</pre></details>


**Any errors, logs, or other information that might help us identify your problem**
