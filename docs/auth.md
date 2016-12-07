# Authentication

Brizo starts up without requiring any user authentication. This mean that the initial user can create an administrative account freely. Please keep this in mind when setting up a new instance of Brizo.

### Github

Currently, authentication is maintained via Github. Brizo needs to be registered as an application for a Github user's account, and then configured to use the Client Id and Client Secret provided. The first user to attempt a login via Github OAuth will also create the administrator account.
