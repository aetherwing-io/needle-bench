class DataProcessor {
  constructor(eventBus) {
    this.eventBus = eventBus;
    this.processedCount = 0;
  }

  process(data) {
    // BUG: Every call to process() registers a NEW listener on the event bus
    // but never removes it. Over time, memory grows unboundedly as closures
    // (and their captured scope) accumulate.
    this.eventBus.on('data:validated', (validatedData) => {
      // This listener captures `data` in its closure, preventing GC
      this._handleValidated(validatedData, data);
    });

    this.eventBus.on('data:error', (err) => {
      console.error(`Processing error for batch ${this.processedCount}: ${err.message}`);
    });

    // Validate and emit
    const validated = this._validate(data);
    if (validated.errors.length > 0) {
      this.eventBus.emit('data:error', new Error(validated.errors.join(', ')));
      return { success: false, errors: validated.errors };
    }

    this.eventBus.emit('data:validated', validated.data);
    this.processedCount++;

    return {
      success: true,
      processed: this.processedCount,
      result: this._transform(validated.data),
    };
  }

  _validate(data) {
    const errors = [];

    if (!data || typeof data !== 'object') {
      errors.push('Data must be a non-null object');
      return { data: null, errors };
    }

    if (!data.items || !Array.isArray(data.items)) {
      errors.push('Missing or invalid items array');
      return { data: null, errors };
    }

    for (let i = 0; i < data.items.length; i++) {
      const item = data.items[i];
      if (!item.id) errors.push(`Item ${i} missing id`);
      if (!item.value && item.value !== 0) errors.push(`Item ${i} missing value`);
    }

    return { data: data, errors };
  }

  _transform(data) {
    return {
      itemCount: data.items.length,
      total: data.items.reduce((sum, item) => sum + (Number(item.value) || 0), 0),
      ids: data.items.map(item => item.id),
    };
  }

  _handleValidated(validatedData, originalData) {
    // Process validated data — the closure retains originalData in memory
    if (validatedData && validatedData.items) {
      // Intentionally captures originalData — this is the leak vector
      void originalData;
    }
  }
}

module.exports = { DataProcessor };
