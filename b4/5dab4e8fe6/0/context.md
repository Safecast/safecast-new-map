# Session Context

## User Prompts

### Prompt 1

Can you check why the global search on the track upload page does not work?

### Prompt 2

Can you check why Github actioon failed to deploy?

### Prompt 3

Can you do that for me?

### Prompt 4

can we deploy one more time to test?

### Prompt 5

[Request interrupted by user for tool use]

### Prompt 6

Aha.. now the old deploy did work!!

### Prompt 7

<task-notification>
<task-id>b20ef7f</task-id>
<output-file>/tmp/claude-1000/-home-rob-Documents-Safecast-safecast-new-map/tasks/b20ef7f.output</output-file>
<status>completed</status>
<summary>Background command "Check if public key is authorized on server" completed (exit code 0)</summary>
</task-notification>
Read the output file to retrieve the result: /tmp/claude-1000/-home-rob-Documents-Safecast-safecast-new-map/tasks/b20ef7f.output

### Prompt 8

After deployemnt the search locally works fine but on the proction not. Could tha be Cloudfront caching?

### Prompt 9

Can you create an empty commit to check if all works?

### Prompt 10

When I visit simplmap.safecast.org  and login, still the global search does not work.?

### Prompt 11

Now I can not clear the search and on admin  user page the search does not work yet?

### Prompt 12

[Request interrupted by user]

### Prompt 13

That is on the producion server..

### Prompt 14

Cleared the cache, tried incognito window, but same results. ?

### Prompt 15

That is not true.. OI can not clear the results and old entry into the search filed come back all the time. Seems the form data iis cached?

### Prompt 16

Did I hard refresh, but the clear button still does not work on production. Locally all work fine..

### Prompt 17

FYI..

CloudFront generally does not cache data submitted through form
POST requests, as its primary function is to cache GET and HEAD requests to improve performance. However, improper configuration, particularly regarding query strings, can lead to unintended caching behavior where form data appears "cached" or user inputs are ignored. 
Here is a breakdown of how CloudFront interacts with form data:

    POST Requests (Standard Forms): By default, CloudFront does not cache responses to POST re...

### Prompt 18

Continue fixng the issues with search

### Prompt 19

[Request interrupted by user]

### Prompt 20

Seems you are starting for cratch for this chat? Please check before chat we had a fe wminutes ago.

### Prompt 21

The last chat ended with :
Prompt is too long

### Prompt 22

[Request interrupted by user for tool use]

### Prompt 23

Can you first read this:

https://entire.io/Safecast/safecast-new-map/checkpoints/main/37a02146df35

### Prompt 24

Please read this /home/rob/Downloads/fix_ Add comment to adminUploadsHandler and force redeploy · Entire.pdf

### Prompt 25

The search in the admin panle does not work on production server. Locally all works fine. This happened after we added CloudFront. Yesterday.

### Prompt 26

That is fixed now..

### Prompt 27

Another issue with caching it seems.. Locally I can upload mulitple log fines, but on the production sever not. See screenshot ( I am logged in).

### Prompt 28

[Request interrupted by user]

### Prompt 29

I checked a bit more, if I log out and then referech the page, the screen and the login seems not to be cleared. And It looks like I am still looged in.

### Prompt 30

Still isues..when Looged in, I can not upload multiple files..single files work fine..

### Prompt 31

No error about the login/cache etc.

### Prompt 32

[Request interrupted by user]

### Prompt 33

One still works fine. Here screnshot from network for the multiple files upload.

