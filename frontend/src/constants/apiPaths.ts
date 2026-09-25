export const apiPaths = {
  artworks: '/api/v1/artworks',
  exhibitions: '/api/v1/exhibitions',
  artists: '/api/v1/artists',
  interactions: '/api/v1/interactions',
};

export const exhibitionReadinessPath = (id: string): string => `${apiPaths.exhibitions}/${id}/readiness`;
