#!/usr/bin/env bash
set -euo pipefail

# The Kubernetes namespace in which Brigade is running.
namespace="default"

event_provider="simple-event"
event_type="my_event"

# This is github.com/deis/empty-testbed
project_id="brigade-adaf165d35013179fb5d1e1e50e0ed7fc99cf00cae67ce389374d5"
commit="master"


# This is the brigade script to execute
script=$(cat <<EOF
const { events } = require("brigadier");
events.on("my_event", (e) => {
  console.log("The system time is " + e.payload);
});
EOF
)

# Now we will generate a new event evrey 60 seconds.
while : ; do
  # We'll use a UUID instead of a ULID. But if you want a ULID generator, you
  # can grab one here: https://github.com/technosophos/ulid
  uuid="$(uuidgen | tr '[:upper:]' '[:lower:]')"

  # We can use the UUID to make sure we get a unique name
  name="simple-event-$uuid"

  # This will just print the system time for the system running the script.
  payload=$(date)

  cat <<EOF | tee test.yaml | kubectl --namespace ${namespace} create -f -
apiVersion: v1
kind: Secret
metadata:
  name: ${name}
  labels:
    heritage: brigade
    project: ${project_id}
    build: ${uuid}
    commit: ${commit}
    component: build
type: Opaque
data:
  commit: $(echo -n "${commit}" | base64 -w0)
  event_provider: $(echo -n "${event_provider}" | base64 -w0)
  event_type: $(echo -n "${event_type}" | base64 -w0)
  project_id: $(echo -n "${project_id}" | base64 -w0)
  build_id: $(echo -n "${uuid}" | base64 -w0)
  payload: $(echo -n "${payload}" | base64 -w0)
  script: $(echo -n "${script}" | base64 -w0)
EOF
  sleep 60
done

