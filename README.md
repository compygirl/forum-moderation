# Forum-image-upload Web Application
* `compygirl` - Aigerim Yessenbayeva




# Forum Image Upload

In the forum image upload, registered users have the possibility to create a post containing an image as well as text.

- When viewing the post, users and guests should see the image associated with it.

There are several extensions for images like: JPEG, SVG, PNG, GIF, etc. In this project, you have to handle at least JPEG, PNG, and GIF types.

The max size of the images to load should be 20 MB. If there is an attempt to load an image greater than 20 MB, an error message should inform the user that the image is too big.

# Forum Security Project

## HTTPS Implementation

### Encrypted Connection

- Generate an SSL certificate: Your website's identity card.
- Choose: Create self-signed certificates or use Certificate Authorities (CA's).
- Cipher Suites: Investigate and select for secure communication.

## Rate Limiting Implementation

- Ensure: Presence of rate limiting against abuse and potential DDoS attacks.
## Objectives

This project aims to develop a web forum with the following functionalities:

- Facilitate communication between users.
- Allow associating categories to posts.
- Enable liking and disliking of posts and comments.
- Implement a filtering mechanism for posts.
- Implement additionally the security






## SQLite Database

To store forum data (users, posts, comments, etc.), the project will utilize the SQLite database library. SQLite is well-suited as an embedded database software for local/client storage in various applications. It allows database creation and management through queries.

For better database structure and performance, it is recommended to refer to an entity relationship diagram and create a corresponding database.

- Usage of at least one SELECT, one CREATE, and one INSERT query is mandatory.

To learn more about SQLite, refer to the [SQLite page](https://www.sqlite.org/index.html).

## Authentication

Users must be able to register as new users for the forum by providing their credentials. Additionally:

- Creation of a login session to access the forum for adding posts and comments.
- Implementation of cookies to ensure a single active session per user with an expiration date.
- Registration should prompt for email, username, and password (bonus task: encryption of passwords).
- Handling of error responses for existing email and incorrect credentials.

## Communication

For communication within the forum:

- Only registered users can create posts and comments.
- Users can associate posts with one or more categories of their choice.
- Posts and comments will be visible to all users (registered or not).
- Non-registered users can only view posts and comments.

## Likes and Dislikes

Liking or disliking posts and comments:

- Exclusive functionality for registered users.
- Display the number of likes and dislikes to all users.

## Filter

A filtering mechanism enabling users to filter displayed posts by:

- Categories
- Created posts
- Liked posts

Registered users can access the last two options, specifically applying to the logged-in user.

## Docker

Building the Docker
```CMD/Terminal 
docker build -t forum .
```

Run Docker
```CMD/Terminal 
docker run --name=forum -p 8080:8080 --rm -d forum
```

Check Running Container
```CMD/Terminal 
docker ps -a
```

Stop runnin container
```CMD/Terminal 
docker stop forum
```
## Instructions

- Compulsory usage of SQLite.
- Proper handling of website errors, HTTP statuses, and technical errors.
- Code adherence to best practices.
- Recommended inclusion of test files for unit testing.

Please replace the placeholders with project-specific information. This README provides a comprehensive overview of the project's objectives, functionalities, and requirements.


## Usage/Examples
Cloning storage to your host
```CMD/Terminal 
git clone git@github.com:compygirl/forum-moderation.git
```
Go to the downloaded repository:

```CMD/Terminal 
cd forum-security
```
Run a Server:
```CMD/Terminal 
go run cmd/main.go
```

Follow the link on the terminal:
```CMD/Terminal 
Starting server got testing... https://localhost:8080 
```

you can play with the page



## HTTP status codes
* OK (200), if everything went without errors.
* Not Found, if the wrong URL was provided.
* Bad Request, for incorrect requests. for example, the id is out of range.
* Internal Server Error, for unhandled errors.



## Feedback

If you liked our project, we would be grateful if you could add `Star` to the repository.

Alem Student
14.02.2024.