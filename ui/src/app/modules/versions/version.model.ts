import { Container } from './containers/container.model';

export class Version {
  public uuid: string;
  public name: string;
  public replicas: number;
  public slug: string;
  public environment_uuid: string;
  public volumes: Volume[];
  public containers: Container[];

  constructor(props: any = {}) {
    Object.assign(this, props);
  }
}

export class Volume {
    public name: string;
    public type: string;

    constructor(props: any = {}) {
      Object.assign(this, props);
      this.type = this.type || "temp";
    }
}
