# wp-userimporter
CLI to import a file of users into a Wordpress installation.

## Usage

### Create a CSV File with your users:

Create a file with a seperate line for each user to be created in the following manner:

```
    first_name, last_name, username, email, password
    Joe, Doe, joedoe, joe@doe.xyz, supersecret123
    ... second user line
    ... and so on ...
```

`username`, `email`, `password` are required.

A headerline as first line is mandatory!

The header fields must be named as described in the Wordpress API!
See here for details:
https://developer.wordpress.org/rest-api/reference/users/#create-a-user

Every further line must contain one user according to syntax of header line.

### Create an application password in your WP Installation:

1. Go to your wordpress admin (`wp-admin`)

1. Select the user profile of the admin user you want to use

1. scroll down to "Application Passwords"

1. create a new application password

### Call `wpuserimporter`:

`wpuserimporter <csv-file> <wp-url> <username> <password>`

`<csv-file>` : the path to the CSV File you just created (see above)

`<wp-url>` : the URL of your Wordpress Installation, pointing to the root-folder, no trailing 
             slash (e.g.: `https://example.com`)

`<username>` : the admin username (from the profile you used for the application password
(see above))

`<password>` : the application password (see above)


# wp-userdeleter
CLI to delete wordpress users by username

## Usage

### Create a file with users to delete

One username per line, no other content

### Create an application password in your WP Installation:

see above (importer)

### Call `wpuserdeleter`:

`wpuserdeleter <file> <wp-url> <username> <password>`

(details like above (importer))

ATTENTION: Presently the deleter will try to move content of deleted user to user with id `0`. 
If the user didn't create any content it is no problem, otherwise there has to be an user with id `0`.
Will be optimized in further version.