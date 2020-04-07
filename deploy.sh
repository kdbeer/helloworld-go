#!/bin/sh

gcloud builds submit --tag gcr.io/facial-app-270721/helloworld --timeout=1200s

gcloud run deploy --image gcr.io/facial-app-270721/helloworld --platform managed