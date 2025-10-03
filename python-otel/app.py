# Import Prometheus client and OTEL metrics libraries
from prometheus_client import start_http_server
from opentelemetry.exporter.prometheus import PrometheusMetricReader
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.resources import SERVICE_NAME, Resource


# Acquire a tracer
tracer = trace.get_tracer("diceroller.tracer")

app = Flask(__name__)

@app.route("/rolldice")
def roll_dice():
  return str(roll())

def roll():
  # This creates a new span that's the child of the current one
  with tracer.start_as_current_span("roll") as rollspan:
    res = randint(1, 6)
    rollspan.set_attribute("roll.value", res)
    return res

# Service name is required for most backends
resource = Resource(attributes={
  SERVICE_NAME: "rolldice"
})

# Start Prometheus client
start_http_server(port=8080, addr="0.0.0.0")
# Initialize PrometheusMetricReader which pulls metrics from the SDK
# on-demand to respond to scrape requests
reader = PrometheusMetricReader()
provider = MeterProvider(resource=resource, metric_readers=[reader])
metrics.set_meter_provider(provider)
