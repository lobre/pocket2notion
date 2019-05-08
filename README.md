# pocket2notion

Currently, there is no way to export Pocket items into Notion because does not provide any API so far.
This projects aims to give a way to append web pages to Notion using the lately introduced [Web Clipper](https://www.notion.so/Web-Clipper-ba54b19ecaeb466b8070b9e683c5fce1).

As the Notion clipper does not give us a way to provide more information than just a title and a URL, we have no way to migrate tags from Pocket to Notion. In order to still pass them, I decided that tags will be appended to the title with a hashtag before. You can disable this feature with the `-t` flag.

Add the `-h` flag to the `pocket2notion` command to see how to use it.

## Notion clipper package

The `clipper` package provides bindings to reproduce the HTTP request made by the clipper extension to add articles to Notion.

### **Disclaimer**
This package has been written by analysing the Chrome extension XHR call. This service exposed by Notion is not public and so there are no guarantees that parameters won't change.
There are even good chances that it will evoluate as Notion might add new features to the clipper in the future.

I have implemented this package to help me doing the switch from Pocket to Notion. I don't aim to keep it updated and to provide full bindings to the Notion service. But feel free to improve/fix it if it turns broken. 

Currently, the package does not support the "+New links database" option as in the extension popup. You can only add items to an already existing database or page. You will by the way need the blockId of this database/page. You can usually find it in the URL of the page when browsing it from the web version of Notion. As well, you will need the value of an authenticated token that you can find using Chrome Developer Tools while browsing authenticated on notion.so. The package does not either give a way to indicate the URL property that you want to use in your existing database in Notion. It will default to choose the first property of type `link`. 

You can find more information in the `clipper` package's Golang documentation.