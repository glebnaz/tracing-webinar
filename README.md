# Tracing Webinar

### Infra

to run jaeger:


`sudo docker run --rm --name jaeger \
-p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:1.8`