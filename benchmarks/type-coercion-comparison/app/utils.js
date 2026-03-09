/**
 * Utility functions for request handling.
 */

function parseQuery(searchParams) {
  const query = {};

  if (searchParams.has('category')) {
    query.category = searchParams.get('category');
  }

  if (searchParams.has('min_price')) {
    query.minPrice = parseFloat(searchParams.get('min_price'));
  }

  if (searchParams.has('max_price')) {
    query.maxPrice = parseFloat(searchParams.get('max_price'));
  }

  if (searchParams.has('min_rating')) {
    // Read the minimum rating from query string
    query.minRating = searchParams.get('min_rating');
  }

  if (searchParams.has('in_stock')) {
    query.inStock = searchParams.get('in_stock') === 'true';
  }

  return query;
}

function formatCurrency(amount) {
  return `$${amount.toFixed(2)}`;
}

module.exports = { parseQuery, formatCurrency };
