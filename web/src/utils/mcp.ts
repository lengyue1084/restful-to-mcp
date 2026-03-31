const configuredGatewayBase = import.meta.env.VITE_MCP_GATEWAY_BASE_URL?.trim();

function trimTrailingSlash(value: string): string {
  return value.replace(/\/+$/, '');
}

export function getGatewayBaseUrl(): string {
  if (configuredGatewayBase) {
    if (/^https?:\/\//i.test(configuredGatewayBase)) {
      return trimTrailingSlash(configuredGatewayBase);
    }
    if (typeof window !== 'undefined') {
      return trimTrailingSlash(new URL(configuredGatewayBase, window.location.origin).toString());
    }
    return trimTrailingSlash(configuredGatewayBase);
  }

  if (typeof window !== 'undefined') {
    return trimTrailingSlash(new URL('/gateway', window.location.origin).toString());
  }

  return '/gateway';
}

export function buildMcpEndpoint(connectToken: string): string {
  const token = connectToken.trim();
  if (!token) return '';
  return `${getGatewayBaseUrl()}/${token}/mcp`;
}
