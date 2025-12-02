#!/bin/bash

# Port-forward:
# $ kubectl port-forward -n mon prometheus-prometheus-pushgateway-665ccd8fbd-7zt9q 9091

JOB="fake-job"

for i in `seq 0 19`
do
  VAL="$((10+RANDOM%20)).$((RANDOM%10))"
  echo "$(date): $i -> $VAL"
  INSTANCE="inst_${i}"
  cat << __EOF | curl --data-binary @- "http://localhost:9091/metrics/job/${JOB}/instance/${INSTANCE}"
# HELP fake_job_load fake load on job
# TYPE fake_job_load gauge
fake_job_load $VAL
__EOF
done
