# Academic Paper API Analysis for Unified Research Search Tool

## Executive Summary

This document provides a comprehensive analysis of major academic paper APIs to design a unified research search tool. We analyze authentication requirements, query capabilities, rate limits, and metadata fields across six major platforms: arXiv, CORE, PubMed, Semantic Scholar, CrossRef, Europe PMC, DOAJ, IEEE Xplore, and PLOS.

## API Comparison Matrix

| API | Authentication | Rate Limits | Response Format | Full-Text Search | Open Access |
|-----|----------------|-------------|-----------------|------------------|-------------|
| **arXiv** | None required | 3-sec delay recommended | XML (Atom 1.0) | No | Yes |
| **CORE** | Optional (Free/Registered) | 1K-200K tokens/day | JSON | Yes | Yes |
| **PubMed** | Optional API key | 3/sec (10/sec with key) | XML/JSON | No | Partial |
| **Semantic Scholar** | Optional API key | 5K/5min (public) | JSON | No | Partial |
| **CrossRef** | None/Optional (Plus) | ~50/sec public | JSON | No | Mixed |
| **Europe PMC** | None required | Not specified | XML/JSON/DC | Yes | Yes |
| **DOAJ** | None (basic) | Not specified | JSON | Yes | Yes |
| **IEEE Xplore** | Required API key | Not specified | JSON | No | Partial |
| **PLOS** | None required | 10/min, 300/hr, 7200/day | JSON (Solr) | Yes | Yes |

## Detailed API Analysis

### 1. arXiv API

**Base URL**: `http://export.arxiv.org/api/query`

**Key Features**:
- No authentication required
- Boolean search operators (AND, OR, ANDNOT)
- Field-specific searches (au:, ti:, abs:, cat:)
- Version-specific article retrieval
- Maximum 2000 results per request
- 30,000 total results cap

**Metadata Fields**:
- title, authors, abstract, categories
- published/updated dates
- journal references (if available)
- links to PDF and abstract

**Search Example**:
```
search_query=au:Einstein+AND+ti:relativity&max_results=10
```

### 2. CORE API v3

**Base URL**: `https://api.core.ac.uk/v3/`

**Key Features**:
- Free tier with authentication levels
- Complex query syntax with field lookups
- Aggregation capabilities
- Deduplication across repositories
- Full-text search capability

**Rate Limit Tracking**:
- X-RateLimitRemaining header
- X-RateLimit-Retry-After header
- X-RateLimit-Limit header

**Searchable Entities**:
- Works, Outputs, Data Providers, Journals

### 3. PubMed E-utilities

**Base URL**: `https://eutils.ncbi.nlm.nih.gov/entrez/eutils/`

**Key Utilities**:
- ESearch: Text queries returning UIDs
- EFetch: Full record retrieval
- ESummary: Document summaries
- ELink: Related records

**Query Syntax**:
```
term[field] AND/OR/NOT term[field]
```

**Best Practices**:
- Run large jobs during off-peak hours
- Include tool and email parameters
- Use API key for higher rate limits

### 4. Semantic Scholar

**Base URL**: `https://api.semanticscholar.org/`

**Key Features**:
- Paper search and batch retrieval
- Citation graph navigation
- Author profiles
- Recommendations API
- No authentication for basic access

**Endpoints**:
- `/graph/v1/paper/search`
- `/graph/v1/paper/{paper_id}`
- `/graph/v1/author/{author_id}`

### 5. CrossRef

**Base URL**: `https://api.crossref.org/`

**Key Features**:
- Anonymous access by default
- "Polite" pool with email header
- Metadata Plus for premium features
- Extensive metadata coverage
- 2000 row limit per request

**Metadata Coverage**:
- Funding data, licenses, full-text links
- ORCID iDs, abstracts, references
- Crossmark updates

### 6. Europe PMC

**Base URL**: `https://www.ebi.ac.uk/europepmc/webservices/rest/`

**Key Features**:
- No registration required
- 33+ million publications
- Text mining annotations API
- Multiple format support (XML/JSON/DC)
- Extends PubMed with patents and guidelines

### 7. DOAJ

**Key Features**:
- Basic access without authentication
- 21,558 journals indexed
- 11+ million article records
- Widget integration available
- XML/CSV upload for publishers

### 8. IEEE Xplore

**Key Features**:
- API key required
- Non-commercial use only
- Metadata access only
- Dynamic query tool for testing

### 9. PLOS

**Base URL**: `http://api.plos.org/search`

**Key Features**:
- Solr-based search
- Article-level metrics API
- CC0 licensed data
- 100 row limit recommended

## Common Metadata Fields for Normalization

Based on the analysis, these fields should be normalized across all APIs:

### Core Fields (Available in Most APIs)
1. **Identifiers**
   - DOI (Digital Object Identifier)
   - PMID (PubMed ID)
   - PMCID (PubMed Central ID)
   - arXiv ID
   - ISBN/ISSN

2. **Basic Metadata**
   - Title
   - Authors (name, ORCID if available)
   - Abstract
   - Publication date
   - Journal/Conference name

3. **Classification**
   - Subject categories
   - Keywords
   - MeSH terms (biomedical)

4. **Access Information**
   - Open access status
   - License information
   - Full-text URLs
   - PDF availability

### Extended Fields (API-specific)
- Citation count
- Reference list
- Funding information
- Clinical trial numbers
- Patent citations

## Design Considerations for Unified Search Tool

### 1. Parallel Search Architecture

**Strategy**: Implement concurrent API queries with intelligent routing

```go
type SearchRequest struct {
    Query      string
    APIs       []string
    MaxResults int
    Filters    SearchFilters
}

type SearchFilters struct {
    DateRange    DateRange
    OpenAccess   bool
    SubjectAreas []string
    Authors      []string
}
```

### 2. Rate Limit Management

**Approach**: Implement per-API rate limiters with backoff strategies

```go
type RateLimiter struct {
    APIName       string
    RequestsPerSec int
    BurstLimit    int
    BackoffStrategy BackoffStrategy
}
```

**Considerations**:
- Use token bucket algorithm for smooth rate limiting
- Implement exponential backoff for 429 errors
- Track rate limit headers where available
- Priority queue for API requests

### 3. Result Aggregation and Deduplication

**Deduplication Strategy**:
1. Primary key: DOI (most reliable)
2. Fallback: Title + First Author + Year
3. Similarity matching for near-duplicates

**Ranking Algorithm**:
- Source reliability score
- Metadata completeness
- Open access preference
- Recency bias (optional)

### 4. Caching Strategy

**Multi-level Cache**:
1. **Request Cache**: Cache full search results (15-minute TTL)
2. **Paper Cache**: Individual papers (24-hour TTL)
3. **Metadata Cache**: Journal info, author profiles (7-day TTL)

### 5. Error Handling and Resilience

**Strategies**:
- Circuit breaker pattern for API failures
- Fallback to cached results
- Partial result returns with error indicators
- Retry with exponential backoff

### 6. Authentication Management

**Approach**: Centralized credential store with environment-based configuration

```go
type APICredentials struct {
    APIName    string
    APIKey     string
    APISecret  string
    TokenType  string // "header", "query", "bearer"
    HeaderName string // e.g., "X-API-Key", "Crossref-Plus-API-Token"
}
```

## Implementation Recommendations

### 1. API Selection Strategy

For comprehensive coverage, prioritize APIs based on:
- **Open Access**: CORE, Europe PMC, DOAJ, PLOS, arXiv
- **Biomedical**: PubMed, Europe PMC
- **General Academic**: CrossRef, Semantic Scholar
- **Physics/Math/CS**: arXiv
- **Engineering**: IEEE Xplore (if commercial use allowed)

### 2. Search Workflow

1. **Query Parsing**: Extract search terms, filters, field specifications
2. **API Selection**: Choose relevant APIs based on subject area
3. **Query Translation**: Convert to API-specific syntax
4. **Parallel Execution**: Launch concurrent searches with timeouts
5. **Result Aggregation**: Merge, deduplicate, and rank results
6. **Response Formatting**: Normalize to unified schema

### 3. Performance Optimizations

- **Connection Pooling**: Reuse HTTP connections
- **Response Streaming**: Process results as they arrive
- **Pagination Strategy**: Lazy loading for large result sets
- **Compression**: Enable gzip for API requests
- **CDN Integration**: Cache static metadata

### 4. Monitoring and Analytics

Track:
- API response times and error rates
- Rate limit utilization
- Cache hit rates
- Search query patterns
- Result quality metrics

## Conclusion

A unified research search tool should implement:

1. **Flexible API abstraction** supporting diverse authentication and query formats
2. **Intelligent rate limiting** with per-API configurations
3. **Robust deduplication** using multiple identifiers
4. **Multi-level caching** for performance
5. **Graceful degradation** when APIs fail
6. **Comprehensive monitoring** for optimization

The tool should prioritize open access sources while maintaining compatibility with subscription-based services, ensuring maximum coverage of the academic literature landscape.