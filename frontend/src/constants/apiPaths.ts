export const apiPaths = {
  artworks: '/api/v1/artworks',
  exhibitions: '/api/v1/exhibitions',
  exhibitionReadiness: (id: string) => `/api/v1/exhibitions/${id}/readiness`,
  artists: '/api/v1/artists',
  interactions: '/api/v1/interactions',
};
