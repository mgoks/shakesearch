# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.onrender.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is in rough shape. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the app! Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful.

You can approach this with a back-end, front-end, or full-stack focus.

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission.

## Submission

1. Fork this repository and send us a link to your fork after pushing your changes.

3. Render (render.com) hosting, the application deploys cleanly from a public url.
   
   Application is live at https://shakesearch-murat.onrender.com/.
   
4. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

   To improve the application I mainly focused on two main areas, namely the search algorithm and the UI.
   
   To improve the search algorithm, I made the search case-insensitive. I updated the algorithm such that it now finds the title of the work in which the queried string is used and sends those data to the front-end along with context. Context is updated such that the algorithm dynamically finds the beginning of the sentence in which is query is used and that is send to the front-end to give the user better understanding of the context. I limited the search query to only 2 or more alphanumeric characters, whitespace, and apostrophes to prevent empty and single-character queries which caused the app to become unresponsive. I added a new type called `Result` to pass search results to the front-end in an encapsulated fashion.
   
   To improve the UI, I used Tailwind CSS. With this UI framework, I created a user-friendly web page by creating a search bar and by adding an image on top of it. The query is now marked in the search results so that the user can easily see the query within the context. I also added a favicon for the page.
   
   The most crucial further improvement at this point is to tolerate typos and misspellings by expanding the search algorithm to include fuzzy text-search. This can be implemented by using the Levenshtein distance algorithm or using an external library if such one exists.
   
   I added a number of *TODO* comments to some of the source files outlining what needs to be done to improve code quality. Here is a list of these tasks here as well.
    * Move the text marking from the back-end to the front-end. This is currently done in the back-end because doing it in the front-end results in a bug in which some results are marked at wrong positions.
    * Change the return type of `Search` function to the `Result` type in main.go.
    * Marshal results to JSON before encoding them in `handleSearch` function in main.go and parse the JSON input when updating the table in app.js. Once that is done, the for-loop in app.js should be converted to a for-of loop or `forEach` function.
    * The search results table should be updated by using DOM and HTML DOM API in app.js instead of writing HTML directly in strings to improve security.

   In addition to the above, I would also add unit tests to both back and front ends, if I had more time.
