upnode
======

Library for distributed task scheduling, example of go code, not working

Purpose of this repository only to showcase coding style, and testing of a library that is work in progress, 
so it might not even compile at this point as a whole, but some individual tests should be working

Also it is missing settings file with machine specific data, for obvious reasons.

Library itself is using core "time" package to schedule events based on information from postgresql 
database, and sending events to beanstalkd queue server, It works in real time, with low overhead and it has been 
tested in production on several machines. 

Part of library that is not completed should handle said events and process them.