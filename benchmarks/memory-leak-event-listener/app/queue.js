class TaskQueue {
  constructor(eventBus) {
    this.eventBus = eventBus;
    this.queue = [];
    this.processing = false;
  }

  enqueue(task) {
    if (!task || !task.type) {
      throw new Error('Task must have a type field');
    }

    this.queue.push({
      ...task,
      enqueuedAt: Date.now(),
      status: 'pending',
    });

    this.eventBus.emit('task:queued', task);
    this._processNext();
  }

  _processNext() {
    if (this.processing || this.queue.length === 0) return;

    this.processing = true;
    const task = this.queue.shift();

    // Simulate async processing
    setTimeout(() => {
      task.status = 'completed';
      task.completedAt = Date.now();
      this.eventBus.emit('task:completed', task);
      this.processing = false;
      this._processNext();
    }, 10);
  }

  getQueueLength() {
    return this.queue.length;
  }
}

module.exports = { TaskQueue };
