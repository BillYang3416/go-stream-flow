import { Component } from '@angular/core';
import { environment } from 'src/enviroments/environment.dev';
import { ApiPrefix, ApiService } from './core/services/api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  title = 'file-flow-ui';

  lineLoginUrl = environment.apiUrl + '/auth/line-login';

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {}

  logout() {
    this.apiService.get(ApiPrefix.AUTH, 'logout', {}).subscribe();
  }
}
