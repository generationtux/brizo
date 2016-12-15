import { ActivatedRoute }     from '@angular/router'
import { Component }          from '@angular/core'

import { ApplicationService } from './application.service'
import { Pod }                from './pod.component'

@Component({
  selector:    'application',
  templateUrl: './application.component.html',
  providers:   [ApplicationService],
})

export class ApplicationComponent {
  private application: Application

  constructor(private applicationService: ApplicationService, private route: ActivatedRoute) {
    let uuid: string
    route.params.subscribe(
      data => uuid = data['uuid'],
    )
    applicationService.getApplication(uuid).subscribe(
      data => this.application = data,
      err => console.error('There was an error: ' + err),
      () => console.log('application loaded'),
    )
  }
}

export class Application {
    constructor(
        public id:   number,
        public uuid: string,
        public name: string,
        public pods: Array<Pod>
    ){}
}
