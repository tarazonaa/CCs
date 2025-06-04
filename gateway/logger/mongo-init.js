db.createUser({
    user: "kong",
    pwd: "holajorge",
    roles: [{
        role: "readWrite",
        db: "mydatabase"
    }]
});
