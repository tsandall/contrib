package example

# allow "admin" to see all posts
allow = true {
    input.user = "admin"
}

allow = true {
    input.method = "GET"
    input.path = ["posts"]
    allowed[x]
}

allow = true {
    input.method = "GET"
    input.path = ["posts", post_id]
    allowed[x]
    x.id = post_id
}

# return posts authored by input.user
allowed[x] {
    data.posts[x]
    x.author == input.user
}

# # return posts with clearance level greater than 0 and less than equal to 5
# # but no posts from "it"
# allowed[x] {
#     x := data.posts[_]
#     x.clearance <= 5
#     x.clearance > 0
#     x.department != "it"
# }

# # return posts containing the term "OPA" in their message
# allowed[x] {
#     x := data.posts[_]
#     contains(x.message, "OPA")
# }

# # return posts who email address matches the ".org" domain
# allowed[x] {
#     x := data.posts[_]
#     re_match("[a-zA-Z]+@[a-zA-Z]+.org", x.email)
# }

# # return posts liked by input.user
# allowed[x] {
#     data.posts[x].likes[_] = input.user
# }