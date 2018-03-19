---
title: "Jiffy Application Structure"
date: 2018-02-10T23:56:10-07:00
weight: 30
draft: true
---

### Jiffy application file structure

This is what Jiffy generates when provided with a model-file describing a simple 'Person' entity.  The structure and files look more or less standard if you are used to looking at such things.  Explanations of each folder and its content are discussed thoughout the documentation.

```code

FirstApp
├── appobj
│   ├── appconf.go
│   └── appobj.go
├── controllers
│   ├── authc.go
│   ├── controllerfuncs.go
│   ├── groupauthc.go
│   ├── person_relationsc.go
│   ├── personc.go
│   ├── usrc.go
│   ├── usr_groupc.go
│   └── ext
│       ├── extc_interfaces.go
│       └── personc_ext.go
├── jwtkeys
│   ├── private.pem
│   └── public.pem
├── middleware
│   └── requireuser.go
├── models
│   ├── authm.go
│   ├── errors.go
│   ├── group_authm.go
│   ├── modelfuncs.go
│   ├── personm_ext.go
│   ├── personm.go
│   ├── servicesm.go
│   ├── usr_groupm.go
│   ├── usrm.go
│   └── ext
│       └── model_ext_interfaces.go
├── util
│   └── strings.go
├── .dev.config.json
├── .prd.config.json
├── main_test.go
└── main.go

```

### Jiffy application services

Jiffy approaches the API from a services perspective.  Each entity has a corresponding service, that can be started when the application initializes.  The Usr, UsrGroup, Auth and GroupAuth services are always generated be default when creating a Jiffy application.  Additional services are generated based on the model files created based on the business scenario.




{{<mermaid align="left">}}
graph LR;
    A[Hard edge] -->|Link text| B(Round edge)
    B --> C{Decision}
    C -->|One| D[Result one]
    C -->|Two| E[Result two]
{{< /mermaid >}}
<br/>

{{<mermaid align="right">}}
graph LR;
    A[Hard edge] -->|Link text| B(Round edge)
    B --> C{Decision}
    C -->|One| D[Result one]
    C -->|Two| E[Result two]
{{< /mermaid >}}
<br/>

{{<mermaid align="center">}}
    graph TD
A[Christmas] -->|Get money| B(Go shopping)
B --> C{Let me think}
C -->|One| D[Laptop]
C -->|Two| E[iPhone]
C -->|Three| F[Car]
{{< /mermaid >}}
<br/>




