/**
 * Product filtering module.
 * Filters products based on query parameters.
 */

function filterProducts(products, query) {
  let results = [...products];

  if (query.category) {
    results = results.filter(p => p.category === query.category);
  }

  if (query.minPrice !== undefined) {
    results = results.filter(p => p.price >= query.minPrice);
  }

  if (query.maxPrice !== undefined) {
    results = results.filter(p => p.price <= query.maxPrice);
  }

  if (query.minRating !== undefined) {
    // Build list of acceptable ratings (minRating through 5)
    const acceptable = buildRatingRange(query.minRating, 5);
    results = results.filter(p => acceptable.includes(p.rating));
  }

  if (query.inStock !== undefined) {
    results = results.filter(p => p.inStock === query.inStock);
  }

  return results;
}

/**
 * Build an array of integers from min to max (inclusive).
 */
function buildRatingRange(min, max) {
  const range = [];
  for (let i = min; i <= max; i++) {
    range.push(i);
  }
  return range;
}

module.exports = { filterProducts, buildRatingRange };
