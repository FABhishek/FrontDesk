# Admin front end setup:
1. npm install
2. npm run dev

# Client front end setup: 
(This code is part of Livekit's documentation, can be ignored from reviewing)
1. pnpm install
2. npm run dev

# Python Microservice 
(used python, as SDK was only available in nodejs and python)
refer to livekit's documentation

# Go backend 
(for API calls, used it as there was language flexibility of choosing any)
1. go build
2. go run.

Database used is supabase(postgres) as I was more familiar with it and was easy to setup.


Note: During screen recording it crashed unexpectedly, although I have covered the most of the important information in video itself, few points I wasn't able to mention:
<img width="1620" height="478" alt="image" src="https://github.com/user-attachments/assets/4857cfce-2264-42e4-b0d2-b9369162cf3d" />
1) We have customer Id available in admin panel, and can record each query associated with customer, so that we can fire webhook easily and notify customers on spot.
2) We can mark the queries as Unresolved using a cron job, which will mark them as unresolved after a certain period of time (assume cron job runs every night)
3) Database schema (could've been more organized, but this was sufficient for the task)
<img width="850" height="794" alt="image" src="https://github.com/user-attachments/assets/90c9219c-72e6-47bc-a3da-1031d04fd9c1" />

