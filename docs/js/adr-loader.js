// ADR Loader - Handles loading ADR files and parsing metadata
window.ADRLoader = (function() {
    'use strict';

    /**
     * Load list of ADRs from various sources
     */
    async function loadADRList() {
        try {
            // Try to load from auto-generated JSON index first (most reliable)
            const jsonAdrs = await loadFromJsonIndex();
            if (jsonAdrs.length > 0) {
                console.info(`Loaded ${jsonAdrs.length} ADRs from JSON index`);
                return jsonAdrs;
            }

            // Fallback: try to load from README
            console.warn('JSON index not found, falling back to README parsing');
            const readmeAdrs = await loadFromReadme();
            if (readmeAdrs.length > 0) {
                console.info(`Loaded ${readmeAdrs.length} ADRs from README`);
                return readmeAdrs;
            }

            // Last resort: try to discover ADRs by scanning directory
            console.warn('README parsing failed, attempting directory discovery');
            const discoveredAdrs = await discoverADRs();
            if (discoveredAdrs.length > 0) {
                console.info(`Discovered ${discoveredAdrs.length} ADRs via fallback method`);
            }
            return discoveredAdrs;

        } catch (error) {
            console.error('Failed to load ADR list:', error);
            throw new Error('Could not load Architecture Decision Records');
        }
    }

    /**
     * Load ADR list from auto-generated JSON index
     */
    async function loadFromJsonIndex() {
        try {
            const response = await fetch('adr-index.json');
            if (!response.ok) {
                throw new Error('JSON index not found');
            }

            const data = await response.json();
            
            if (!data.adrs || !Array.isArray(data.adrs)) {
                throw new Error('Invalid JSON index format');
            }

            console.info(`JSON index generated: ${data.generated}`);
            return data.adrs.sort((a, b) => a.number.localeCompare(b.number));

        } catch (error) {
            console.warn('Could not load from JSON index:', error);
            return [];
        }
    }

    /**
     * Load ADR list from README.md
     */
    async function loadFromReadme() {
        try {
            const response = await fetch('../README.md');
            if (!response.ok) {
                throw new Error('README not found');
            }

            const readmeContent = await response.text();
            return parseReadmeIndex(readmeContent);

        } catch (error) {
            console.warn('Could not load from README:', error);
            return [];
        }
    }

    /**
     * Parse ADR index from README content
     */
    function parseReadmeIndex(readmeContent) {
        const adrs = [];
        
        // Look for the ADR index table
        const lines = readmeContent.split('\n');
        let inTable = false;
        let headerPassed = false;

        for (const line of lines) {
            // Start of table
            if (line.includes('| ADR | Title | Status |')) {
                inTable = true;
                continue;
            }

            // Header separator
            if (inTable && line.includes('|-----')) {
                headerPassed = true;
                continue;
            }

            // End of table
            if (inTable && headerPassed && (!line.trim().startsWith('|') || line.includes('###'))) {
                break;
            }

            // Parse table row
            if (inTable && headerPassed && line.trim().startsWith('|')) {
                const adr = parseReadmeTableRow(line);
                if (adr) {
                    adrs.push(adr);
                }
            }
        }

        return adrs.sort((a, b) => a.number.localeCompare(b.number));
    }

    /**
     * Parse individual table row from README
     */
    function parseReadmeTableRow(line) {
        const cells = line.split('|').map(cell => cell.trim()).filter(cell => cell);
        
        if (cells.length < 3) return null;

        try {
            // Extract ADR number and file path from the first cell
            const adrCell = cells[0];
            const linkMatch = adrCell.match(/\[(\d{4})\]\(([^)]+)\)/);
            
            if (!linkMatch) return null;

            const number = linkMatch[1];
            const filePath = linkMatch[2];
            const title = cells[1];
            const status = cells[2];
            const diagramType = cells[3] || '-';

            // Convert relative path to be relative to docs/ directory
            // README paths are like "docs/adr/0001-file.md" but we need "adr/0001-file.md"
            const relativePath = filePath.startsWith('docs/adr/') ? 
                filePath.replace('docs/', '') : 
                filePath.startsWith('adr/') ? 
                filePath : 
                `adr/${filePath.split('/').pop()}`;

            return {
                number,
                title,
                status,
                diagramType,
                filePath: relativePath,
                fileName: filePath.split('/').pop()
            };

        } catch (error) {
            console.warn('Failed to parse README row:', line, error);
            return null;
        }
    }

    /**
     * Discover ADRs by probing for files systematically
     * Since we can't list directory contents in browsers, we'll try a different approach:
     * attempt to fetch files using various common naming patterns
     */
    async function discoverADRs() {
        const adrs = [];
        const adrPattern = /^(\d{4})-(.+)\.md$/;

        console.info('Attempting ADR discovery fallback - this may take a moment...');

        // Common ADR filename patterns/prefixes observed in the wild
        const commonPatterns = [
            'record-architecture-decisions',
            'establish-architecture-review-board', 
            'adopt-microservices-architecture',
            'choose-database-per-service',
            'implement-api-gateway-pattern',
            'use-event-driven-communication',
            'implement-graphql-api',
            'use-mongodb-for-session-storage',
            'use-redis-for-session-storage',
            'adopt-hybrid-session-storage',
            // Common verbs + generic patterns for unknown ADRs
            'adopt-', 'use-', 'implement-', 'choose-', 'establish-', 'record-',
            'select-', 'define-', 'migrate-', 'replace-', 'remove-', 'add-',
            'update-', 'change-', 'deprecate-', 'retire-', 'introduce-'
        ];

        // Try sequential numbers up to 20 (reasonable limit)
        for (let i = 1; i <= 20; i++) {
            const number = i.toString().padStart(4, '0');
            let foundFile = false;

            // First try exact patterns we know about
            for (const pattern of commonPatterns) {
                if (foundFile) break;
                
                const filePath = pattern.endsWith('-') ? 
                    null : // Skip partial patterns for first pass
                    `adr/${number}-${pattern}.md`;
                
                if (!filePath) continue;

                try {
                    const response = await fetch(filePath);
                    if (response.ok) {
                        const content = await response.text();
                        const metadata = parseADRMetadata(content);
                        const fileName = filePath.split('/').pop();
                        
                        adrs.push({
                            number,
                            title: metadata.title || generateTitleFromFilename(fileName),
                            status: metadata.status || 'Unknown',
                            diagramType: metadata.diagramType || '-',
                            filePath,
                            fileName
                        });
                        
                        foundFile = true;
                        console.info(`Found ADR ${number}: ${metadata.title}`);
                        break;
                    }
                } catch (e) {
                    // File doesn't exist, continue
                    continue;
                }
            }
        }

        if (adrs.length === 0) {
            console.warn('No ADRs discovered via fallback method. README.md parsing is recommended.');
        } else {
            console.info(`Discovered ${adrs.length} ADRs via fallback method`);
        }

        return adrs.sort((a, b) => a.number.localeCompare(b.number));
    }

    /**
     * Generate a title from filename when metadata parsing fails
     */
    function generateTitleFromFilename(fileName) {
        const match = fileName.match(/^\d{4}-(.+)\.md$/);
        if (!match) return fileName;
        
        return match[1]
            .split('-')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    }

    /**
     * Load content of a specific ADR
     */
    async function loadADRContent(adr) {
        try {
            const response = await fetch(adr.filePath);
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const content = await response.text();
            
            if (!content.trim()) {
                throw new Error('ADR file is empty');
            }

            return content;

        } catch (error) {
            console.error(`Failed to load ADR content for ${adr.number}:`, error);
            throw new Error(`Could not load ADR ${adr.number}: ${error.message}`);
        }
    }

    /**
     * Parse metadata from ADR content
     */
    function parseADRMetadata(content) {
        const metadata = {
            title: null,
            status: null,
            diagramType: '-'
        };

        const lines = content.split('\n');

        for (let i = 0; i < lines.length; i++) {
            const line = lines[i].trim();

            // Extract title (first h1)
            if (!metadata.title && line.startsWith('# ')) {
                metadata.title = line.substring(2).trim();
            }

            // Extract status
            if (line === '## Status' && i + 1 < lines.length) {
                // Look for status in next few lines
                for (let j = i + 1; j < Math.min(i + 5, lines.length); j++) {
                    const statusLine = lines[j].trim();
                    if (statusLine && !statusLine.startsWith('#') && !statusLine.startsWith('##')) {
                        metadata.status = statusLine;
                        break;
                    }
                }
            }

            // Detect diagram types
            if (line.includes('```mermaid') || line.includes('mermaid')) {
                if (line.includes('C4Context') || content.includes('C4Context')) {
                    metadata.diagramType = 'Context';
                } else if (line.includes('C4Container') || content.includes('C4Container')) {
                    metadata.diagramType = 'Container';
                } else if (line.includes('C4Component') || content.includes('C4Component')) {
                    metadata.diagramType = 'Component';
                } else if (line.includes('sequenceDiagram') || content.includes('sequenceDiagram') || 
                          line.includes('flowchart') || content.includes('flowchart')) {
                    metadata.diagramType = 'Dynamic';
                }
            }
        }

        return metadata;
    }

    /**
     * Validate ADR content structure
     */
    function validateADRStructure(content) {
        const requiredSections = ['## Status', '## Context', '## Decision', '## Consequences'];
        const missingsections = [];

        for (const section of requiredSections) {
            if (!content.includes(section)) {
                missingSection.push(section);
            }
        }

        return {
            isValid: missingSections.length === 0,
            missingSection: missingSection
        };
    }

    /**
     * Get ADR statistics
     */
    function getADRStatistics(adrs) {
        const stats = {
            total: adrs.length,
            byStatus: {},
            byDiagramType: {},
            withDiagrams: 0
        };

        adrs.forEach(adr => {
            // Count by status
            const status = adr.status || 'Unknown';
            stats.byStatus[status] = (stats.byStatus[status] || 0) + 1;

            // Count by diagram type
            const diagramType = adr.diagramType || '-';
            stats.byDiagramType[diagramType] = (stats.byDiagramType[diagramType] || 0) + 1;

            // Count ADRs with diagrams
            if (diagramType !== '-') {
                stats.withDiagrams++;
            }
        });

        return stats;
    }

    /**
     * Search ADRs by content
     */
    async function searchADRContent(adrs, searchTerm) {
        const results = [];
        const term = searchTerm.toLowerCase();

        for (const adr of adrs) {
            try {
                const content = await loadADRContent(adr);
                const contentLower = content.toLowerCase();

                if (contentLower.includes(term)) {
                    // Find context around matches
                    const matches = [];
                    const lines = content.split('\n');
                    
                    lines.forEach((line, index) => {
                        if (line.toLowerCase().includes(term)) {
                            matches.push({
                                line: index + 1,
                                text: line.trim(),
                                context: lines.slice(Math.max(0, index - 1), index + 2)
                            });
                        }
                    });

                    results.push({
                        adr,
                        matches
                    });
                }
            } catch (error) {
                console.warn(`Could not search content of ADR ${adr.number}:`, error);
            }
        }

        return results;
    }

    /**
     * Get related ADRs based on content similarity
     */
    function getRelatedADRs(currentAdr, allAdrs, maxResults = 3) {
        if (!currentAdr.content) return [];

        const currentWords = extractKeywords(currentAdr.content);
        const scored = [];

        allAdrs.forEach(adr => {
            if (adr.number === currentAdr.number) return;

            let score = 0;

            // Score based on title similarity
            const titleWords = adr.title.toLowerCase().split(/\W+/);
            titleWords.forEach(word => {
                if (currentWords.includes(word) && word.length > 3) {
                    score += 3;
                }
            });

            // Score based on status
            if (adr.status === currentAdr.status) {
                score += 1;
            }

            // Score based on diagram type
            if (adr.diagramType === currentAdr.diagramType && adr.diagramType !== '-') {
                score += 2;
            }

            if (score > 0) {
                scored.push({ adr, score });
            }
        });

        return scored
            .sort((a, b) => b.score - a.score)
            .slice(0, maxResults)
            .map(item => item.adr);
    }

    /**
     * Extract keywords from content
     */
    function extractKeywords(content) {
        const text = content.toLowerCase()
            .replace(/[^\w\s]/g, ' ')
            .replace(/\s+/g, ' ');
        
        const words = text.split(' ')
            .filter(word => word.length > 3)
            .filter(word => !isStopWord(word));

        // Count word frequency
        const frequency = {};
        words.forEach(word => {
            frequency[word] = (frequency[word] || 0) + 1;
        });

        // Return most frequent words
        return Object.keys(frequency)
            .sort((a, b) => frequency[b] - frequency[a])
            .slice(0, 20);
    }

    /**
     * Check if word is a stop word
     */
    function isStopWord(word) {
        const stopWords = new Set([
            'this', 'that', 'with', 'have', 'will', 'been', 'from', 'they', 'know',
            'want', 'been', 'good', 'much', 'some', 'time', 'very', 'when', 'come',
            'here', 'just', 'like', 'long', 'make', 'many', 'over', 'such', 'take',
            'than', 'them', 'well', 'were', 'what', 'your', 'section', 'decision',
            'architecture', 'system', 'service', 'application', 'implementation',
            'approach', 'solution', 'requirements', 'design', 'technology'
        ]);
        
        return stopWords.has(word);
    }

    // Public API
    return {
        loadADRList,
        loadADRContent,
        parseADRMetadata,
        validateADRStructure,
        getADRStatistics,
        searchADRContent,
        getRelatedADRs
    };

})();