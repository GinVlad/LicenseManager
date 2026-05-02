import SwaggerUI from 'swagger-ui-react';
import 'swagger-ui-react/swagger-ui.css';

export default function ApiDocs() {
  return (
    <div className="swagger-wrapper">
      <SwaggerUI url="/swagger.json" />
    </div>
  );
}
