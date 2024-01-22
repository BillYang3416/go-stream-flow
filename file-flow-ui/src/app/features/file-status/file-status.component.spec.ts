import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FileStatusComponent } from './file-status.component';

describe('FileStatusComponent', () => {
  let component: FileStatusComponent;
  let fixture: ComponentFixture<FileStatusComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FileStatusComponent]
    });
    fixture = TestBed.createComponent(FileStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
