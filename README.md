# Feature-toggle

Small application that simulates the Feature toggle app.

## Disclaimers

- I could've used a real DB solution instead of implementing the DB functionality, but I intentionally opted for the 
  map solution with an interface that can later be satisfied by the DB adapter (PostgreSQL, MongoDB etc);
- I also could've put the FE app behind nginx and proxy requests to the BE, but it would've taken more time, so I
  decided to serve static files through the BE server (CORS isn't cool, but for a test exercise will do).

## HOWTO

The full control over the app is available through the [Makefile](Makefile).
In order to run the app in Docker, you can either use docker-compose or run it yourself (users are advised to use former).

The UI app provides the full CRUD interface for the features.
To create a feature, use the form. In order to update a feature, press the 'edit' button of a feature from a table.

#### IMPORTANT:
Although a feature can be created with only technical name and customer ids fields, its updating requires four fields 
to be filled: description name, expiration date, description and feature status.
I realise that it sounds like a bad design; had I more time, I would've done things differently.
