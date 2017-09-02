
  // ##### GET, SET STATE

  // ## get original message from the queue (lpop the list)

  // ## Save the message to a key with name queue-id.task-id
  //    SET queue-id.task-id <task>

  // ## Save the (queue-id.task-id) to a REDIS Set working_set with
  //    sort date now + 5 min
  //    ZADD Q_working_set <now>+300 queue-id.task-id    

  // ## send it to the client


  // ##### UPDATE / HEARTBEAT

  // ## Receive queue-id and task-id from client
  // ## Update the (queue-id.task-id) with a new expires
  //    ZADD Q_working_set <now>+300 queue-id.task-id


  // ##### COMPLETE / FINISH

  // ## Receive queue-id and task-id from client "Done"
  // ## Remove the index from the working_set
  //    ZREM Q_working_set "queue-id.task-id"

  // ## Add it to the list "done"
  //    ZADD Q_completed_set task-id queue-id.task-id


  // ##### GET FAILED
  // ## this will actually be done FIRST every time a consumer wants to get a
  //    new item from the queue.

  // ## Get all items from working set that have expired
  //    ZRANGEBYSCORE Q_working_set -inf <now>

  // ## Increment the counter on the item
  //    INCR queue-id.task-id.count

  // If the counter < max max tries:
  // Return the item to the consumer.
  // Update the timeout.

  
