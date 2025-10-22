# Admin front end setup:
1. npm install
2. npm run dev

# Client front end setup:
1. pnpm install
2. npm run dev

# Python Microservice
refer to livekit's documentation

# Go backend (for API calls)
1. go build
2. go run.

Database used is supabase(postgres) as I was more familiar with it and was easy to setup.


Note: During screen recording it crashed unexpectedly, although I have covered the most of the important information in video itself, few points I wasn't able to mention:
<img width="1620" height="478" alt="image" src="https://github.com/user-attachments/assets/4857cfce-2264-42e4-b0d2-b9369162cf3d" />
1) We have customer Id available in admin panel, and can record each query associated with customer, so that we can fire webhook easily and notify customers on spot.
2) We can mark the queries as Unresolved using a cron job, which will mark them as unresolved after a certain period of time (assume cron job runs every night)
