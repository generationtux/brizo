import { ActivatedRoute, Params }     from '@angular/router'
import { Component, OnInit }          from '@angular/core'
import 'rxjs/add/operator/switchMap';

import { ApplicationService } from './application.service'
import { Pod }                from '../pod.component'

@Component({
  selector:    'application',
  templateUrl: './application-details.component.html',
  providers:   [ApplicationService],
})

export class ApplicationComponent implements OnInit {
  private application: Application

  constructor(private applicationService: ApplicationService, private route: ActivatedRoute) {}

  ngOnInit() {
    this.route.params
      .switchMap(
        (params: Params) => this.applicationService.getApplication(params['uuid'])
      ).subscribe(
        data => this.application = data,
        err => console.error('There was an error: ' + err)
      );
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
