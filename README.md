# DevOps Recruitment Test 

## Instruction
In this recruitment test you will use a IaC (Infrastructure as a Code) to solve realworld problems in system architecture.
The goals use of IaC in this test is to have a reproducible and has more consistent environment.

Requirements :
- Reproducible environment, so your result can be evaluated and measured.
- Deploy the sample application with provided architecture.
- Setup monitoring tools.
- Setup distributed tracing.
- Configure alert for the collected metrics via email.
- Able to scale up and down the sample application properly.

### Docker compose 
Docker Compose is used for the IaC in this test, which provide simple setup even on local machine.

### Stacks 
These are the required stack to be deployed :
 - redis: in-memory data structure store.
 - prometheus: metrics store and collector
 - grafana: metrics visualization
 - jaeger: distributed tracing
 - nginx: web server and loadbalancer
 <!-- - loki:  -->

 You can add additional stacks that you wanna use, but please include the clarification

### App Architecture
These are application that need to be deployed :
- counter-service: contain a service that provide a counter

Following diagaram are the overview architecture of the application that are need to be setup
![App architecture](/docs/images/SampleApp%20-%20Architecture.drawio.svg)

### Monitoring and Alert goals
The goals of monitoring are we need to ensure that the application are running well after deployment.
In such case in this test there are some measurement that being enabled, such as :
- To make sure that the web server/loadbalancer
- Achieve service uptime SLA by monitor from the user experience
- Monitor host resource usage
- Application services are healthy
- System dependency such as DB are also healthy and under expected performance

In such case to further emphasize our goals we also need to introduce alerting systems whenever :
- Application services are not healthy
- DB perfomance degraded
- Unable to fulfill the expected SLA

### Tracing
And to help investigation in case performance degradation or failure we need to know which parts of the system that has been impacted