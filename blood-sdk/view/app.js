/* app.js */

// require and instantiate express
const app = require('express')()
var nets = require("nets")
// fake posts to simulate a database

// set the view engine to ejs
app.set('view engine', 'ejs');
app.set('views',__dirname+'/');

let posts_data = {
    cp_name_1 : "radhe",
    cp_name_2 : "vidit",
    cp_name_3 : "kapil",
    cp_name_4 : "dhruva",
   };
let hos1 = [
        {
            grp : 'A+',
            count : 20,
            requires : 21

        },
        {
            grp : 'B',
            count : 20,
            requires : 20
        },
        {
            grp : 'AB+',
            count : 50,
            requires : 20
        },
        {
            grp : 'O+',
            count : 50,
            requires : 20
        }
    ];
   let hos2 = [
    {
        grp : 'A+',
        count : 20,
        requires : 21

    },
    {
        grp : 'B',
        count : 20,
        requires : 20
    },
    {
        grp : 'AB+',
        count : 50,
        requires : 20
    },
    {
        grp : 'O+',
        count : 50,
        requires : 20
    }
    ];
   let hos3 = [
    {
        grp : 'A+',
        count : 20,
        requires : 21

    },
    {
        grp : 'B',
        count : 20,
        requires : 20
    },
    {
        grp : 'AB+',
        count : 50,
        requires : 20
    },
    {
        grp : 'O+',
        count : 50,
        requires : 20
    }
    ];
  let  hos4 = [
    {
        grp : 'A+',
        count : 20,
        requires : 21

    },
    {
        grp : 'B',
        count : 20,
        requires : 20
    },
    {
        grp : 'AB+',
        count : 50,
        requires : 20
    },
    {
        grp : 'O+',
        count : 50,
        requires : 20
    }
    ];
    
 
// blog home page
app.get('/', (req, res) => {
  // render `home.ejs` with the list of posts
  nets({ url: "http://placekitten.com/g/400/400" }, function(err, resp, body) {
  // body is a Buffer containing the image
})
  res.render('home', { posts: posts_data, hos_1 : hos1, hos_2 : hos2, hos_3 : hos3,hos_4 : hos4 })
})



// blog post
app.get('/post/:id', (req, res) => {
  // find the post in the `posts` array
  const post = posts.filter((post) => {
    return post.id == req.params.id
  });

  // render the `post.ejs` template with the post content
  res.render('post', {
    author: post.author,
    title: post.title,
    body: post.body
  })
})

app.listen(8080)

console.log('listening on port 8080')